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
	"github.com/nadeeshame/repograph_platform/pkg/utils"
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
	defer logger.Sync()

	logger.Info("Starting Document Scanner Service",
		zap.String("version", "1.0.0"),
		zap.Int("port", 8081))

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"healthy": true})
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/scan/directory", scanDirectory)
		v1.GET("/metadata/:filePath", getFileMetadata)
		v1.POST("/compute-hash", computeHash)
	}

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		logger.Info("Document Scanner service starting", zap.String("address", srv.Addr))
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

func scanDirectory(c *gin.Context) {
	var req struct {
		Directory string `json:"directory" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Scanning directory", zap.String("directory", req.Directory))

	// TODO: Implement directory scanning
	c.JSON(http.StatusOK, gin.H{
		"directory": req.Directory,
		"files":     []string{},
	})
}

func getFileMetadata(c *gin.Context) {
	filePath := c.Param("filePath")
	logger.Info("Getting file metadata", zap.String("file_path", filePath))

	// TODO: Implement metadata extraction
	c.JSON(http.StatusOK, gin.H{
		"path": filePath,
		"type": utils.GetFileCategory(filePath),
	})
}

func computeHash(c *gin.Context) {
	var req struct {
		FilePath string `json:"file_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := utils.ComputeFileHash(req.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_path": req.FilePath,
		"hash":      hash,
	})
}
