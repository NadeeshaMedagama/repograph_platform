package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Azure    AzureConfig    `mapstructure:"azure"`
	Google   GoogleConfig   `mapstructure:"google"`
	Pinecone PineconeConfig `mapstructure:"pinecone"`
	GitHub   GitHubConfig   `mapstructure:"github"`
	App      AppConfig      `mapstructure:"app"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Services ServicesConfig `mapstructure:"services"`
	Server   ServerConfig   `mapstructure:"server"`
}

// AzureConfig contains Azure OpenAI configuration
type AzureConfig struct {
	OpenAIAPIKey               string `mapstructure:"openai_api_key"`
	OpenAIEndpoint             string `mapstructure:"openai_endpoint"`
	OpenAIEmbeddingsVersion    string `mapstructure:"openai_embeddings_version"`
	OpenAIEmbeddingsDeployment string `mapstructure:"openai_embeddings_deployment"`
	OpenAIAPIVersion           string `mapstructure:"openai_api_version"`
	OpenAIChatDeployment       string `mapstructure:"openai_chat_deployment"`
}

// GoogleConfig contains Google Vision API configuration
type GoogleConfig struct {
	VisionAPIKey           string `mapstructure:"vision_api_key"`
	ApplicationCredentials string `mapstructure:"application_credentials"`
}

// PineconeConfig contains Pinecone vector database configuration
type PineconeConfig struct {
	APIKey        string `mapstructure:"api_key"`
	Host          string `mapstructure:"host"`
	IndexName     string `mapstructure:"index_name"`
	Dimension     int    `mapstructure:"dimension"`
	Cloud         string `mapstructure:"cloud"`
	Region        string `mapstructure:"region"`
	UseNamespaces bool   `mapstructure:"use_namespaces"`
}

// GitHubConfig contains GitHub API configuration
type GitHubConfig struct {
	Token string `mapstructure:"token"`
}

// AppConfig contains application-level configuration
type AppConfig struct {
	DataDirectory         string `mapstructure:"data_directory"`
	LogLevel              string `mapstructure:"log_level"`
	ChunkSize             int    `mapstructure:"chunk_size"`
	ChunkOverlap          int    `mapstructure:"chunk_overlap"`
	SkipExistingDocuments bool   `mapstructure:"skip_existing_documents"`
}

// RedisConfig contains Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// ServicesConfig contains microservices URLs
type ServicesConfig struct {
	DocumentScannerURL      string `mapstructure:"document_scanner_url"`
	ContentExtractorURL     string `mapstructure:"content_extractor_url"`
	VisionServiceURL        string `mapstructure:"vision_service_url"`
	SummarizationServiceURL string `mapstructure:"summarization_service_url"`
	EmbeddingServiceURL     string `mapstructure:"embedding_service_url"`
	VectorStoreServiceURL   string `mapstructure:"vector_store_service_url"`
	QueryServiceURL         string `mapstructure:"query_service_url"`
	OrchestratorServiceURL  string `mapstructure:"orchestrator_service_url"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// Load loads configuration from environment and config files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Set defaults
	setDefaults()

	// Automatically read environment variables
	viper.AutomaticEnv()

	// Bind environment variables with proper naming
	bindEnvVariables()

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// Azure OpenAI defaults
	viper.SetDefault("azure.openai_embeddings_version", "2024-02-01")
	viper.SetDefault("azure.openai_api_version", "2024-02-01")
	viper.SetDefault("azure.openai_embeddings_deployment", "text-embedding-ada-002")
	viper.SetDefault("azure.openai_chat_deployment", "gpt-4")

	// Pinecone defaults
	viper.SetDefault("pinecone.dimension", 1536)
	viper.SetDefault("pinecone.cloud", "aws")
	viper.SetDefault("pinecone.region", "us-east-1")
	viper.SetDefault("pinecone.use_namespaces", true)

	// Application defaults
	viper.SetDefault("app.data_directory", "./data/diagrams")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.chunk_size", 1000)
	viper.SetDefault("app.chunk_overlap", 200)
	viper.SetDefault("app.skip_existing_documents", true)

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 30*time.Second)
	viper.SetDefault("server.write_timeout", 30*time.Second)

	// Service URLs defaults
	viper.SetDefault("services.document_scanner_url", "http://localhost:8081")
	viper.SetDefault("services.content_extractor_url", "http://localhost:8082")
	viper.SetDefault("services.vision_service_url", "http://localhost:8083")
	viper.SetDefault("services.summarization_service_url", "http://localhost:8084")
	viper.SetDefault("services.embedding_service_url", "http://localhost:8085")
	viper.SetDefault("services.vector_store_service_url", "http://localhost:8086")
	viper.SetDefault("services.query_service_url", "http://localhost:8087")
	viper.SetDefault("services.orchestrator_service_url", "http://localhost:8088")
}

func bindEnvVariables() {
	// viper.BindEnv binds environment variables to configuration keys
	// These bindings are used when reading config values

	// Azure OpenAI
	viper.BindEnv("azure.openai_api_key", "AZURE_OPENAI_API_KEY")                             //nolint:errcheck
	viper.BindEnv("azure.openai_endpoint", "AZURE_OPENAI_ENDPOINT")                           //nolint:errcheck
	viper.BindEnv("azure.openai_embeddings_version", "AZURE_OPENAI_EMBEDDINGS_VERSION")       //nolint:errcheck
	viper.BindEnv("azure.openai_embeddings_deployment", "AZURE_OPENAI_EMBEDDINGS_DEPLOYMENT") //nolint:errcheck
	viper.BindEnv("azure.openai_api_version", "AZURE_OPENAI_API_VERSION")                     //nolint:errcheck
	viper.BindEnv("azure.openai_chat_deployment", "AZURE_OPENAI_CHAT_DEPLOYMENT")             //nolint:errcheck

	// Google
	viper.BindEnv("google.vision_api_key", "GOOGLE_VISION_API_KEY")                   //nolint:errcheck
	viper.BindEnv("google.application_credentials", "GOOGLE_APPLICATION_CREDENTIALS") //nolint:errcheck

	// Pinecone
	viper.BindEnv("pinecone.api_key", "PINECONE_API_KEY")               //nolint:errcheck
	viper.BindEnv("pinecone.host", "PINECONE_HOST")                     //nolint:errcheck
	viper.BindEnv("pinecone.index_name", "PINECONE_INDEX_NAME")         //nolint:errcheck
	viper.BindEnv("pinecone.dimension", "PINECONE_DIMENSION")           //nolint:errcheck
	viper.BindEnv("pinecone.cloud", "PINECONE_CLOUD")                   //nolint:errcheck
	viper.BindEnv("pinecone.region", "PINECONE_REGION")                 //nolint:errcheck
	viper.BindEnv("pinecone.use_namespaces", "PINECONE_USE_NAMESPACES") //nolint:errcheck

	// GitHub
	viper.BindEnv("github.token", "GITHUB_TOKEN") //nolint:errcheck

	// App
	viper.BindEnv("app.data_directory", "DATA_DIRECTORY")                   //nolint:errcheck
	viper.BindEnv("app.log_level", "LOG_LEVEL")                             //nolint:errcheck
	viper.BindEnv("app.chunk_size", "CHUNK_SIZE")                           //nolint:errcheck
	viper.BindEnv("app.chunk_overlap", "CHUNK_OVERLAP")                     //nolint:errcheck
	viper.BindEnv("app.skip_existing_documents", "SKIP_EXISTING_DOCUMENTS") //nolint:errcheck

	// Redis
	viper.BindEnv("redis.host", "REDIS_HOST")         //nolint:errcheck
	viper.BindEnv("redis.port", "REDIS_PORT")         //nolint:errcheck
	viper.BindEnv("redis.password", "REDIS_PASSWORD") //nolint:errcheck

	// Server
	viper.BindEnv("server.port", "SERVICE_PORT") //nolint:errcheck

	// Services
	viper.BindEnv("services.document_scanner_url", "DOCUMENT_SCANNER_URL")           //nolint:errcheck
	viper.BindEnv("services.content_extractor_url", "CONTENT_EXTRACTOR_URL")         //nolint:errcheck
	viper.BindEnv("services.vision_service_url", "VISION_SERVICE_URL")               //nolint:errcheck
	viper.BindEnv("services.summarization_service_url", "SUMMARIZATION_SERVICE_URL") //nolint:errcheck
	viper.BindEnv("services.embedding_service_url", "EMBEDDING_SERVICE_URL")         //nolint:errcheck
	viper.BindEnv("services.vector_store_service_url", "VECTOR_STORE_SERVICE_URL")   //nolint:errcheck
	viper.BindEnv("services.query_service_url", "QUERY_SERVICE_URL")                 //nolint:errcheck
	viper.BindEnv("services.orchestrator_service_url", "ORCHESTRATOR_SERVICE_URL")   //nolint:errcheck
}

func validate(config *Config) error {
	// Required: Azure OpenAI configuration
	if config.Azure.OpenAIAPIKey == "" {
		return fmt.Errorf("AZURE_OPENAI_API_KEY is required")
	}
	if config.Azure.OpenAIEndpoint == "" {
		return fmt.Errorf("AZURE_OPENAI_ENDPOINT is required")
	}

	// Required: Pinecone configuration
	if config.Pinecone.APIKey == "" {
		return fmt.Errorf("PINECONE_API_KEY is required")
	}
	if config.Pinecone.IndexName == "" {
		return fmt.Errorf("PINECONE_INDEX_NAME is required")
	}

	// Required: Application configuration
	if config.App.DataDirectory == "" {
		return fmt.Errorf("DATA_DIRECTORY is required")
	}

	// Validate numeric constraints
	if config.App.ChunkSize <= 0 {
		return fmt.Errorf("chunk_size must be positive")
	}
	if config.App.ChunkOverlap < 0 {
		return fmt.Errorf("chunk_overlap cannot be negative")
	}
	if config.App.ChunkOverlap >= config.App.ChunkSize {
		return fmt.Errorf("chunk_overlap must be less than chunk_size")
	}
	if config.Pinecone.Dimension <= 0 {
		return fmt.Errorf("pinecone dimension must be positive")
	}

	// Note: Google Vision API key is optional
	// Note: GitHub token is optional

	return nil
}

// GetRedisAddr returns the Redis connection address
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
