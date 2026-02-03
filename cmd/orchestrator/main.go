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
	"github.com/nadeeshame/repograph_platform/internal/orchestrator"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Initialize logger
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer func() { _ = logger.Sync() }() //nolint:errcheck

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", zap.Error(err))
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	logger.Info("Starting Orchestrator Service",
		zap.String("version", "1.0.0"),
		zap.Int("port", cfg.Server.Port))

	// Setup HTTP router
	router := gin.Default()

	// Health endpoint - simple check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "orchestrator",
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Ready endpoint
	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ready": true})
	})

	// API endpoints
	v1 := router.Group("/api/v1")
	{
		v1.POST("/process/document", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "processing"})
		})
		v1.POST("/process/directory", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "processing"})
		})
		v1.GET("/status/:documentId", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "unknown"})
		})
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Server starting", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Start automatic indexing in background
	if cfg.App.DataDirectory != "" {
		go func() {
			// Wait for services to be ready
			logger.Info("Waiting for services to be ready before indexing...")
			time.Sleep(10 * time.Second)

			logger.Info("Starting automatic document indexing",
				zap.String("directory", cfg.App.DataDirectory),
				zap.Bool("skip_existing", cfg.App.SkipExistingDocuments))

			// Create document processor
			processor, procErr := orchestrator.NewDocumentProcessor(cfg, logger)
			if procErr != nil {
				logger.Error("Failed to create document processor", zap.Error(procErr))
				return
			}

			// Process directory
			ctx := context.Background()
			procErr = processor.ProcessDirectory(ctx, cfg.App.DataDirectory)
			if procErr != nil {
				logger.Error("Failed to process directory", zap.Error(procErr))
			} else {
				logger.Info("Automatic indexing completed successfully")
			}
		}()
	} else {
		logger.Warn("DATA_DIRECTORY not set, automatic indexing disabled")
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("Server exited")
	return nil
}
