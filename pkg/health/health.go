package health

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// Status represents the health status of a service
type Status struct {
	Healthy   bool              `json:"healthy"`
	Services  map[string]bool   `json:"services"`
	Timestamp time.Time         `json:"timestamp"`
	Details   map[string]string `json:"details,omitempty"`
}

// Checker provides health checking functionality
type Checker struct {
	azureEndpoint   string
	pineconeAPIKey  string
	googleVisionKey string
	db              *sql.DB
	redisClient     *redis.Client
}

// NewChecker creates a new health checker
func NewChecker(azureEndpoint, pineconeAPIKey, googleVisionKey string, db *sql.DB, redisClient *redis.Client) *Checker {
	return &Checker{
		azureEndpoint:   azureEndpoint,
		pineconeAPIKey:  pineconeAPIKey,
		googleVisionKey: googleVisionKey,
		db:              db,
		redisClient:     redisClient,
	}
}

// CheckAll checks the health of all services
func (c *Checker) CheckAll(ctx context.Context) *Status {
	services := make(map[string]bool)
	details := make(map[string]string)

	// Check Azure OpenAI
	services["azure_openai"] = c.checkAzureOpenAI(ctx)
	if !services["azure_openai"] {
		details["azure_openai"] = "Unable to reach Azure OpenAI endpoint"
	}

	// Check Pinecone
	services["pinecone"] = c.checkPinecone(ctx)
	if !services["pinecone"] {
		details["pinecone"] = "Unable to verify Pinecone connection"
	}

	// Check Google Vision
	services["google_vision"] = c.checkGoogleVision(ctx)
	if !services["google_vision"] {
		details["google_vision"] = "Unable to verify Google Vision API"
	}

	// Check Database
	services["database"] = c.checkDatabase(ctx)
	if !services["database"] {
		details["database"] = "Unable to connect to database"
	}

	// Check Redis
	services["redis"] = c.checkRedis(ctx)
	if !services["redis"] {
		details["redis"] = "Unable to connect to Redis"
	}

	// Overall health is true if all services are healthy
	healthy := true
	for _, status := range services {
		if !status {
			healthy = false
			break
		}
	}

	return &Status{
		Healthy:   healthy,
		Services:  services,
		Timestamp: time.Now(),
		Details:   details,
	}
}

func (c *Checker) checkAzureOpenAI(ctx context.Context) bool {
	if c.azureEndpoint == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Simple HTTP check to the endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.azureEndpoint, nil)
	if err != nil {
		return false
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Azure OpenAI may return 401 if no auth, but the endpoint is reachable
	return resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusOK
}

func (c *Checker) checkPinecone(_ context.Context) bool {
	// For now, just check if API key is configured
	// In a real implementation, you'd make an API call to Pinecone
	return c.pineconeAPIKey != ""
}

func (c *Checker) checkGoogleVision(_ context.Context) bool {
	// For now, just check if API key is configured
	// In a real implementation, you'd make an API call to Google Vision
	return c.googleVisionKey != ""
}

func (c *Checker) checkDatabase(ctx context.Context) bool {
	if c.db == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := c.db.PingContext(ctx); err != nil {
		return false
	}
	return true
}

func (c *Checker) checkRedis(ctx context.Context) bool {
	if c.redisClient == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := c.redisClient.Ping(ctx).Err(); err != nil {
		return false
	}
	return true
}

// CheckAzureOpenAI checks Azure OpenAI service
func (c *Checker) CheckAzureOpenAI(ctx context.Context) bool {
	return c.checkAzureOpenAI(ctx)
}

// CheckPinecone checks Pinecone service
func (c *Checker) CheckPinecone(ctx context.Context) bool {
	return c.checkPinecone(ctx)
}

// CheckGoogleVision checks Google Vision service
func (c *Checker) CheckGoogleVision(ctx context.Context) bool {
	return c.checkGoogleVision(ctx)
}

// CheckDatabase checks database connection
func (c *Checker) CheckDatabase(ctx context.Context) bool {
	return c.checkDatabase(ctx)
}

// CheckRedis checks Redis connection
func (c *Checker) CheckRedis(ctx context.Context) bool {
	return c.checkRedis(ctx)
}

// HTTPHandler returns an HTTP handler for health checks
func (c *Checker) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := c.CheckAll(r.Context())

		statusCode := http.StatusOK
		if !status.Healthy {
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		// Write JSON response
		fmt.Fprintf(w, `{"healthy":%v,"timestamp":"%s","services":%v}`,
			status.Healthy,
			status.Timestamp.Format(time.RFC3339),
			marshalServices(status.Services))
	}
}

func marshalServices(services map[string]bool) string {
	result := "{"
	first := true
	for k, v := range services {
		if !first {
			result += ","
		}
		result += fmt.Sprintf(`"%s":%v`, k, v)
		first = false
	}
	result += "}"
	return result
}
