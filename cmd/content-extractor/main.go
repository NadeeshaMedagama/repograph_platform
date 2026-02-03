package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nadeeshame/repograph_platform/internal/config"
	"github.com/nadeeshame/repograph_platform/internal/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	if err := logger.Initialize(cfg.App.LogLevel); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = logger.Sync() }()

	logger.Info("Starting Content Extractor Service",
		zap.String("version", "1.0.0"),
		zap.Int("port", 8082))

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"healthy": true})
	})

	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ready": true})
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/extract", extractContent)
		v1.GET("/formats", getSupportedFormats)
	}

	srv := &http.Server{
		Addr:         ":8082",
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		logger.Info("Content Extractor service starting", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func extractContent(c *gin.Context) {
	var req struct {
		FilePath string                 `json:"file_path" binding:"required"`
		FileType string                 `json:"file_type" binding:"required"`
		Options  map[string]interface{} `json:"options"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Extracting content",
		zap.String("file_path", req.FilePath),
		zap.String("file_type", req.FileType))

	// TODO: Implement content extraction logic
	c.JSON(http.StatusOK, gin.H{
		"content":      "Extracted content placeholder",
		"extracted_at": time.Now(),
	})
}

func getSupportedFormats(c *gin.Context) {
	// TODO: Return actual supported formats from processors
	formats := []map[string]interface{}{
		{
			"category":   "image",
			"extensions": []string{"png", "jpg", "jpeg", "svg", "gif"},
		},
		{
			"category":   "document",
			"extensions": []string{"pdf", "docx", "pptx", "txt", "md"},
		},
		{
			"category":   "spreadsheet",
			"extensions": []string{"xlsx", "xls", "csv"},
		},
		{
			"category":   "code",
			"extensions": []string{"go", "py", "js", "ts", "java"},
		},
	}

	c.JSON(http.StatusOK, gin.H{"formats": formats})
}
