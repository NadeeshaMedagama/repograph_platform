package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/nadeeshame/rag-knowledge-service/internal/domain/models"
)

// DocumentScanner is responsible for scanning and discovering files
type DocumentScanner interface {
	ScanDirectory(ctx context.Context, directory string) ([]*models.FileMetadata, error)
	GetFileMetadata(ctx context.Context, filePath string) (*models.FileMetadata, error)
	ComputeFileHash(ctx context.Context, filePath string) (string, error)
}

// ContentExtractor extracts content from various file formats
type ContentExtractor interface {
	Extract(ctx context.Context, filePath string, fileType string) (string, error)
	SupportsFileType(fileType string) bool
	GetSupportedFormats() []string
}

// VisionAnalyzer analyzes images and visual content using Google Vision API
type VisionAnalyzer interface {
	AnalyzeImage(ctx context.Context, imageData []byte, imagePath string) (string, error)
	AnalyzeDiagram(ctx context.Context, imageData []byte) (string, error)
	DetectText(ctx context.Context, imageData []byte) (string, error)
}

// SummarizationService generates summaries using Azure OpenAI
type SummarizationService interface {
	Summarize(ctx context.Context, content string, maxTokens int) (string, error)
	SummarizeWithContext(ctx context.Context, content string, context string, maxTokens int) (string, error)
}

// EmbeddingService generates embeddings using Azure OpenAI
type EmbeddingService interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
	GenerateEmbeddings(ctx context.Context, texts []string) ([][]float32, error)
	GetDimension() int
}

// ChunkerService splits text into chunks for embedding
type ChunkerService interface {
	ChunkText(ctx context.Context, text string, chunkSize, overlap int) ([]*models.Chunk, error)
	ChunkDocument(ctx context.Context, doc *models.Document, chunkSize, overlap int) ([]*models.Chunk, error)
}

// VectorStore manages vector storage and retrieval in Pinecone
type VectorStore interface {
	UpsertVectors(ctx context.Context, chunks []*models.Chunk) error
	SearchSimilar(ctx context.Context, queryEmbedding []float32, topK int, namespace string, filter map[string]interface{}) ([]*models.SearchResult, error)
	DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error
	GetIndexStats(ctx context.Context) (map[string]interface{}, error)
	CheckDocumentExists(ctx context.Context, documentHash string) (bool, error)
}

// QueryService handles RAG queries
type QueryService interface {
	Query(ctx context.Context, query *models.Query) (*models.QueryResult, error)
	SearchDocuments(ctx context.Context, query *models.Query) ([]*models.SearchResult, error)
}

// DocumentRepository manages document persistence
type DocumentRepository interface {
	Create(ctx context.Context, doc *models.Document) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Document, error)
	GetByHash(ctx context.Context, hash string) (*models.Document, error)
	Update(ctx context.Context, doc *models.Document) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Document, error)
	GetByState(ctx context.Context, state models.ProcessingState) ([]*models.Document, error)
}

// ChunkRepository manages chunk persistence
type ChunkRepository interface {
	Create(ctx context.Context, chunk *models.Chunk) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Chunk, error)
	GetByDocumentID(ctx context.Context, documentID uuid.UUID) ([]*models.Chunk, error)
	Update(ctx context.Context, chunk *models.Chunk) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByDocumentID(ctx context.Context, documentID uuid.UUID) error
}

// Orchestrator coordinates the document processing workflow
type Orchestrator interface {
	ProcessDocument(ctx context.Context, filePath string) error
	ProcessDirectory(ctx context.Context, directory string, forceReprocess bool) error
	GetProcessingStatus(ctx context.Context, documentID uuid.UUID) (*models.Document, error)
}

// HealthChecker checks the health of external services
type HealthChecker interface {
	CheckHealth(ctx context.Context) map[string]bool
	CheckAzureOpenAI(ctx context.Context) bool
	CheckPinecone(ctx context.Context) bool
	CheckGoogleVision(ctx context.Context) bool
	CheckDatabase(ctx context.Context) bool
	CheckRedis(ctx context.Context) bool
}
