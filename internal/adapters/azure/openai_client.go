package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nadeeshame/repograph_platform/internal/config"
	"go.uber.org/zap"
)

// OpenAIClient handles Azure OpenAI operations
type OpenAIClient struct {
	apiKey              string
	endpoint            string
	embeddingDeployment string
	chatDeployment      string
	apiVersion          string
	httpClient          *http.Client
	logger              *zap.Logger
}

// EmbeddingRequest represents the request body for embeddings
type EmbeddingRequest struct {
	Input []string `json:"input"`
}

// EmbeddingResponse represents the response from embeddings API
type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

// ChatRequest represents the request body for chat completions
type ChatRequest struct {
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float32       `json:"temperature,omitempty"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents the response from chat API
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewOpenAIClient creates a new Azure OpenAI client
func NewOpenAIClient(cfg *config.Config, logger *zap.Logger) (*OpenAIClient, error) {
	if cfg.Azure.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("azure OpenAI API key is required")
	}
	if cfg.Azure.OpenAIEndpoint == "" {
		return nil, fmt.Errorf("azure OpenAI endpoint is required")
	}

	return &OpenAIClient{
		apiKey:              cfg.Azure.OpenAIAPIKey,
		endpoint:            cfg.Azure.OpenAIEndpoint,
		embeddingDeployment: cfg.Azure.OpenAIEmbeddingsDeployment,
		chatDeployment:      cfg.Azure.OpenAIChatDeployment,
		apiVersion:          cfg.Azure.OpenAIAPIVersion,
		httpClient:          &http.Client{},
		logger:              logger,
	}, nil
}

// GenerateEmbedding creates embeddings for the given text
func (c *OpenAIClient) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	c.logger.Debug("Generating embedding", zap.Int("text_length", len(text)))

	url := fmt.Sprintf("%s/openai/deployments/%s/embeddings?api-version=%s",
		c.endpoint, c.embeddingDeployment, c.apiVersion)

	reqBody := EmbeddingRequest{Input: []string{text}}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) //nolint:errcheck
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var embResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(embResp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	c.logger.Debug("Embedding generated successfully",
		zap.Int("dimensions", len(embResp.Data[0].Embedding)))

	return embResp.Data[0].Embedding, nil
}

// GenerateSummary generates a summary for the given text
func (c *OpenAIClient) GenerateSummary(ctx context.Context, text string) (string, error) {
	if text == "" {
		return "", fmt.Errorf("text cannot be empty")
	}

	// Truncate text if too long
	if len(text) > 10000 {
		text = text[:10000] + "..."
	}

	c.logger.Debug("Generating summary", zap.Int("text_length", len(text)))

	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s",
		c.endpoint, c.chatDeployment, c.apiVersion)

	reqBody := ChatRequest{
		Messages: []ChatMessage{
			{Role: "system", Content: "You are a helpful assistant that creates concise, informative summaries. Focus on key points and main ideas."},
			{Role: "user", Content: fmt.Sprintf("Please provide a comprehensive summary of the following content:\n\n%s", text)},
		},
		MaxTokens:   500,
		Temperature: 0.3,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) //nolint:errcheck
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no summary generated")
	}

	summary := chatResp.Choices[0].Message.Content
	c.logger.Debug("Summary generated successfully", zap.Int("summary_length", len(summary)))

	return summary, nil
}

// ChatCompletion performs a chat completion
func (c *OpenAIClient) ChatCompletion(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userMessage},
	}

	reqBody := ChatRequest{
		Messages:    messages,
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s",
		c.endpoint, c.chatDeployment, c.apiVersion)

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) //nolint:errcheck
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return chatResp.Choices[0].Message.Content, nil
}
