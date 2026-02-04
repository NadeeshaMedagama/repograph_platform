package pinecone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/nadeeshame/rag-knowledge-service/internal/config"
	"go.uber.org/zap"
)

// PineconeClient handles Pinecone vector database operations
type PineconeClient struct {
	apiKey     string
	host       string
	httpClient *http.Client
	config     *config.PineconeConfig
	logger     *zap.Logger
}

// Vector represents a vector with metadata
type Vector struct {
	ID       string                 `json:"id"`
	Values   []float32              `json:"values"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Match represents a search result
type Match struct {
	ID       string                 `json:"id"`
	Score    float32                `json:"score"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// UpsertRequest represents the upsert request body
type UpsertRequest struct {
	Vectors   []*Vector `json:"vectors"`
	Namespace string    `json:"namespace,omitempty"`
}

// UpsertResponse represents the upsert response
type UpsertResponse struct {
	UpsertedCount int `json:"upsertedCount"`
}

// QueryRequest represents the query request body
type QueryRequest struct {
	Vector          []float32              `json:"vector"`
	TopK            int                    `json:"topK"`
	IncludeMetadata bool                   `json:"includeMetadata"`
	Filter          map[string]interface{} `json:"filter,omitempty"`
	Namespace       string                 `json:"namespace,omitempty"`
}

// QueryResponse represents the query response
type QueryResponse struct {
	Matches []*QueryMatch `json:"matches"`
}

// QueryMatch represents a match from query
type QueryMatch struct {
	ID       string                 `json:"id"`
	Score    float32                `json:"score"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewPineconeClient creates a new Pinecone client
func NewPineconeClient(cfg *config.Config, logger *zap.Logger) (*PineconeClient, error) {
	if cfg.Pinecone.APIKey == "" {
		return nil, fmt.Errorf("pinecone API key is required")
	}
	if cfg.Pinecone.IndexName == "" {
		return nil, fmt.Errorf("pinecone index name is required")
	}

	httpClient := &http.Client{}
	var host string

	// Use provided host or fetch from Pinecone API
	if cfg.Pinecone.Host != "" {
		host = cfg.Pinecone.Host
		if !strings.HasPrefix(host, "https://") {
			host = "https://" + host
		}
		logger.Info("Using provided Pinecone host", zap.String("host", host))
	} else {
		// Fetch host from Pinecone control plane API
		fetchedHost, err := fetchIndexHost(context.Background(), httpClient, cfg.Pinecone.APIKey, cfg.Pinecone.IndexName, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Pinecone index host: %w", err)
		}
		host = fetchedHost
	}

	logger.Info("Connecting to Pinecone",
		zap.String("index", cfg.Pinecone.IndexName),
		zap.String("host", host))

	return &PineconeClient{
		apiKey:     cfg.Pinecone.APIKey,
		host:       host,
		httpClient: httpClient,
		config:     &cfg.Pinecone,
		logger:     logger,
	}, nil
}

// fetchIndexHost fetches the actual host URL from Pinecone control plane API
func fetchIndexHost(ctx context.Context, httpClient *http.Client, apiKey, indexName string, logger *zap.Logger) (string, error) {
	url := fmt.Sprintf("https://api.pinecone.io/indexes/%s", indexName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Accept", "application/json")

	logger.Debug("Fetching Pinecone index info", zap.String("url", url))

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("API error (status %d): failed to read response body: %w", resp.StatusCode, err)
		}
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var indexInfo struct {
		Host string `json:"host"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&indexInfo); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if indexInfo.Host == "" {
		return "", fmt.Errorf("no host returned from Pinecone API")
	}

	host := indexInfo.Host
	if !strings.HasPrefix(host, "https://") {
		host = "https://" + host
	}

	logger.Info("Fetched Pinecone index host", zap.String("host", host))
	return host, nil
}

// UpsertVectors upserts multiple vectors to Pinecone
func (c *PineconeClient) UpsertVectors(ctx context.Context, vectors []*Vector) error {
	if len(vectors) == 0 {
		return nil
	}

	c.logger.Debug("Upserting vectors", zap.Int("count", len(vectors)))

	// Upsert in batches of 100
	batchSize := 100
	for i := 0; i < len(vectors); i += batchSize {
		end := i + batchSize
		if end > len(vectors) {
			end = len(vectors)
		}

		batch := vectors[i:end]
		if err := c.upsertBatch(ctx, batch); err != nil {
			return fmt.Errorf("failed to upsert batch: %w", err)
		}

		c.logger.Debug("Upserted batch",
			zap.Int("from", i),
			zap.Int("to", end))
	}

	c.logger.Info("Successfully upserted vectors",
		zap.Int("total", len(vectors)))

	return nil
}

func (c *PineconeClient) upsertBatch(ctx context.Context, vectors []*Vector) error {
	reqBody := UpsertRequest{Vectors: vectors}
	if c.config.UseNamespaces {
		reqBody.Namespace = "default"
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/vectors/upsert", c.host)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("API error (status %d): failed to read response body: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// QueryVectors searches for similar vectors
func (c *PineconeClient) QueryVectors(ctx context.Context, embedding []float32, topK int, filter map[string]interface{}) ([]*Match, error) {
	c.logger.Debug("Querying vectors",
		zap.Int("topK", topK),
		zap.Bool("has_filter", filter != nil))

	reqBody := QueryRequest{
		Vector:          embedding,
		TopK:            topK,
		IncludeMetadata: true,
		Filter:          filter,
	}
	if c.config.UseNamespaces {
		reqBody.Namespace = "default"
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/query", c.host)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) //nolint:errcheck
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var queryResp QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&queryResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to Match slice
	matches := make([]*Match, len(queryResp.Matches))
	for i, m := range queryResp.Matches {
		matches[i] = &Match{
			ID:       m.ID,
			Score:    m.Score,
			Metadata: m.Metadata,
		}
	}

	c.logger.Debug("Query complete", zap.Int("matches", len(matches)))
	return matches, nil
}

// CheckDocumentExists checks if a document with given hash exists
func (c *PineconeClient) CheckDocumentExists(ctx context.Context, fileHash string) (bool, error) {
	c.logger.Debug("Checking document existence", zap.String("hash", fileHash))

	// Create a dummy vector for querying
	dummyVector := make([]float32, c.config.Dimension)

	filter := map[string]interface{}{
		"file_hash": map[string]interface{}{
			"$eq": fileHash,
		},
	}

	matches, err := c.QueryVectors(ctx, dummyVector, 1, filter)
	if err != nil {
		// If filter query fails, return false (assume not exists)
		c.logger.Warn("Failed to check document existence", zap.Error(err))
		return false, nil
	}

	exists := len(matches) > 0
	c.logger.Debug("Document existence check complete",
		zap.Bool("exists", exists))

	return exists, nil
}

// GetStats returns index statistics
func (c *PineconeClient) GetStats(ctx context.Context) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/describe_index_stats", c.host)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) //nolint:errcheck
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return stats, nil
}
