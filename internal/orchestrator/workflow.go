package orchestrator

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nadeeshame/rag-knowledge-service/internal/adapters/azure"
	"github.com/nadeeshame/rag-knowledge-service/internal/adapters/google"
	"github.com/nadeeshame/rag-knowledge-service/internal/adapters/pinecone"
	"github.com/nadeeshame/rag-knowledge-service/internal/config"
	"github.com/nadeeshame/rag-knowledge-service/internal/content-extractor/processors"
	"go.uber.org/zap"
)

// DocumentProcessor handles the complete document processing workflow
type DocumentProcessor struct {
	azureClient    *azure.OpenAIClient
	visionClient   *google.VisionClient
	pineconeClient *pinecone.PineconeClient
	processors     []processors.ProcessorInterface
	config         *config.Config
	logger         *zap.Logger
}

// NewDocumentProcessor creates a new document processor
func NewDocumentProcessor(
	cfg *config.Config,
	logger *zap.Logger,
) (*DocumentProcessor, error) {
	// Initialize Azure OpenAI client
	azureClient, err := azure.NewOpenAIClient(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure client: %w", err)
	}

	// Initialize Google Vision client (optional)
	var visionClient *google.VisionClient
	if cfg.Google.VisionAPIKey != "" || cfg.Google.ApplicationCredentials != "" {
		visionClient, err = google.NewVisionClient(cfg, logger)
		if err != nil {
			logger.Warn("Failed to create Vision client", zap.Error(err))
		}
	}

	// Initialize Pinecone client
	pineconeClient, err := pinecone.NewPineconeClient(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Pinecone client: %w", err)
	}

	// Initialize content processors
	contentProcessors := []processors.ProcessorInterface{
		processors.NewTextProcessor(logger),
		processors.NewImageProcessor(logger),
		processors.NewDocumentProcessor(logger),
		processors.NewSpreadsheetProcessor(logger),
		processors.NewCodeProcessor(logger),
	}

	return &DocumentProcessor{
		azureClient:    azureClient,
		visionClient:   visionClient,
		pineconeClient: pineconeClient,
		processors:     contentProcessors,
		config:         cfg,
		logger:         logger,
	}, nil
}

// ProcessDirectory processes all files in a directory
func (dp *DocumentProcessor) ProcessDirectory(ctx context.Context, directory string) error {
	dp.logger.Info("Starting directory processing", zap.String("directory", directory))

	// Scan directory
	files, err := dp.scanDirectory(directory)
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	dp.logger.Info("Found files", zap.Int("count", len(files)))

	// Process each file
	successCount := 0
	skipCount := 0
	errorCount := 0

	for i, file := range files {
		dp.logger.Info("Processing file",
			zap.Int("index", i+1),
			zap.Int("total", len(files)),
			zap.String("file", file))

		err := dp.processFile(ctx, file)
		if err != nil {
			if strings.Contains(err.Error(), "already indexed") {
				skipCount++
				dp.logger.Info("Skipped already-indexed file", zap.String("file", file))
			} else {
				errorCount++
				dp.logger.Error("Failed to process file",
					zap.String("file", file),
					zap.Error(err))
			}
			continue
		}

		successCount++
	}

	dp.logger.Info("Directory processing complete",
		zap.Int("total_files", len(files)),
		zap.Int("processed", successCount),
		zap.Int("skipped", skipCount),
		zap.Int("errors", errorCount))

	return nil
}

// scanDirectory recursively scans a directory for files
func (dp *DocumentProcessor) scanDirectory(directory string) ([]string, error) {
	var files []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and hidden files
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}

// processFile processes a single file
//
//nolint:gocyclo
func (dp *DocumentProcessor) processFile(ctx context.Context, filePath string) error {
	// Calculate file hash
	fileHash, err := dp.calculateFileHash(filePath)
	if err != nil {
		return fmt.Errorf("failed to calculate hash: %w", err)
	}

	// Check if already indexed
	if dp.config.App.SkipExistingDocuments {
		exists, existsErr := dp.pineconeClient.CheckDocumentExists(ctx, fileHash)
		if existsErr != nil {
			dp.logger.Warn("Failed to check document existence", zap.Error(existsErr))
		} else if exists {
			return fmt.Errorf("file already indexed")
		}
	}

	// Extract content
	content, err := dp.extractContent(ctx, filePath)
	if err != nil {
		return fmt.Errorf("failed to extract content: %w", err)
	}

	if content == "" {
		dp.logger.Warn("No content extracted", zap.String("file", filePath))
		return nil
	}

	// Analyze image if applicable
	visualContent := ""
	if dp.isImageFile(filePath) && dp.visionClient != nil {
		var visionErr error
		visualContent, visionErr = dp.visionClient.AnalyzeImage(ctx, filePath)
		if visionErr != nil {
			dp.logger.Warn("Failed to analyze image", zap.Error(visionErr))
		}
	}

	// Combine content
	combinedContent := content
	if visualContent != "" {
		combinedContent += "\n\n" + visualContent
	}

	// Generate summary
	summary, err := dp.azureClient.GenerateSummary(ctx, combinedContent)
	if err != nil {
		dp.logger.Warn("Failed to generate summary", zap.Error(err))
		summary = "Summary generation failed"
	}

	// Create chunks
	chunks := dp.chunkText(combinedContent, dp.config.App.ChunkSize, dp.config.App.ChunkOverlap)

	// Generate document ID
	docID := uuid.New().String()

	// Process each chunk
	vectors := make([]*pinecone.Vector, 0, len(chunks))
	for i, chunk := range chunks {
		// Generate embedding
		chunkEmbedding, embErr := dp.azureClient.GenerateEmbedding(ctx, chunk)
		if embErr != nil {
			dp.logger.Error("Failed to generate embedding",
				zap.Int("chunk", i),
				zap.Error(embErr))
			continue
		}

		// Create vector
		vectorID := fmt.Sprintf("%s-chunk-%d", docID, i)
		vector := &pinecone.Vector{
			ID:     vectorID,
			Values: chunkEmbedding,
			Metadata: map[string]interface{}{
				"document_id": docID,
				"file_name":   filepath.Base(filePath),
				"file_path":   filePath,
				"file_type":   filepath.Ext(filePath),
				"file_hash":   fileHash,
				"chunk_index": i,
				"chunk_total": len(chunks),
				"content":     chunk,
				"summary":     summary,
				"indexed_at":  time.Now().Unix(),
			},
		}

		vectors = append(vectors, vector)
	}

	// Store in Pinecone
	if len(vectors) > 0 {
		err = dp.pineconeClient.UpsertVectors(ctx, vectors)
		if err != nil {
			return fmt.Errorf("failed to store in Pinecone: %w", err)
		}

		dp.logger.Info("Successfully indexed file",
			zap.String("file", filepath.Base(filePath)),
			zap.Int("chunks", len(vectors)))
	}

	return nil
}

// extractContent extracts content using appropriate processor
func (dp *DocumentProcessor) extractContent(ctx context.Context, filePath string) (string, error) {
	ext := filepath.Ext(filePath)

	for _, processor := range dp.processors {
		if processor.CanProcess(ext) {
			return processor.Extract(ctx, filePath)
		}
	}

	return "", fmt.Errorf("no processor found for file type: %s", ext)
}

// isImageFile checks if file is an image
func (dp *DocumentProcessor) isImageFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	imageExts := []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".svg"}
	for _, imgExt := range imageExts {
		if ext == imgExt {
			return true
		}
	}
	return false
}

// calculateFileHash calculates SHA-256 hash of file
func (dp *DocumentProcessor) calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("sha256:%x", hash.Sum(nil)), nil
}

// chunkText splits text into overlapping chunks
func (dp *DocumentProcessor) chunkText(text string, chunkSize, overlap int) []string {
	if len(text) <= chunkSize {
		return []string{text}
	}

	var chunks []string
	start := 0

	for start < len(text) {
		end := start + chunkSize
		if end > len(text) {
			end = len(text)
		}

		chunks = append(chunks, text[start:end])

		start += chunkSize - overlap
	}

	return chunks
}
