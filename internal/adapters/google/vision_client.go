package google

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nadeeshame/repograph_platform/internal/config"
	"go.uber.org/zap"
)

// VisionClient handles Google Vision API operations
type VisionClient struct {
	apiKey string
	logger *zap.Logger
}

// NewVisionClient creates a new Google Vision client
func NewVisionClient(cfg *config.Config, logger *zap.Logger) (*VisionClient, error) {
	// Vision API is optional
	if cfg.Google.VisionAPIKey == "" {
		logger.Warn("Google Vision API key not configured, image analysis will be limited")
	}

	return &VisionClient{
		apiKey: cfg.Google.VisionAPIKey,
		logger: logger,
	}, nil
}

// AnalyzeImage analyzes an image and returns a description
func (c *VisionClient) AnalyzeImage(ctx context.Context, imagePath string) (string, error) {
	c.logger.Debug("Analyzing image", zap.String("path", imagePath))

	// Read the image file
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	// Get file extension to determine type
	ext := strings.ToLower(filepath.Ext(imagePath))

	// For SVG files, just return the content as text
	if ext == ".svg" {
		return c.parseSVG(imageData), nil
	}

	// If no API key, return basic info
	if c.apiKey == "" {
		return c.getBasicImageInfo(imagePath, imageData), nil
	}

	// Use Vision API (simplified - just return basic analysis for now)
	return c.getBasicImageInfo(imagePath, imageData), nil
}

// DetectText extracts text from an image using OCR
func (c *VisionClient) DetectText(ctx context.Context, imagePath string) (string, error) {
	c.logger.Debug("Detecting text in image", zap.String("path", imagePath))

	// Read the image file
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	// For SVG files, extract text content
	ext := strings.ToLower(filepath.Ext(imagePath))
	if ext == ".svg" {
		return c.extractSVGText(imageData), nil
	}

	// If no API key, return empty
	if c.apiKey == "" {
		c.logger.Warn("Vision API not configured, skipping OCR")
		return "", nil
	}

	// Return empty for now (Vision API would be called here)
	return "", nil
}

// AnalyzeDiagram analyzes a diagram and returns structured information
func (c *VisionClient) AnalyzeDiagram(ctx context.Context, imagePath string) (string, error) {
	c.logger.Debug("Analyzing diagram", zap.String("path", imagePath))

	// Get basic image analysis
	analysis, err := c.AnalyzeImage(ctx, imagePath)
	if err != nil {
		return "", err
	}

	// Get text from image
	text, err := c.DetectText(ctx, imagePath)
	if err != nil {
		c.logger.Warn("Failed to detect text", zap.Error(err))
	}

	if text != "" {
		analysis += "\n\nExtracted Text:\n" + text
	}

	return analysis, nil
}

// Helper functions

func (c *VisionClient) parseSVG(data []byte) string {
	content := string(data)

	// Extract text elements from SVG
	var texts []string
	// Simple text extraction - look for text between tags
	parts := strings.Split(content, ">")
	for _, part := range parts {
		if idx := strings.Index(part, "<"); idx > 0 {
			text := strings.TrimSpace(part[:idx])
			if text != "" && !strings.HasPrefix(text, "<?") && !strings.HasPrefix(text, "<!") {
				texts = append(texts, text)
			}
		}
	}

	if len(texts) > 0 {
		return "SVG Diagram Content:\n" + strings.Join(texts, "\n")
	}
	return "SVG Diagram (no text content extracted)"
}

func (c *VisionClient) extractSVGText(data []byte) string {
	// Simple text extraction from SVG
	content := string(data)
	var texts []string

	// Look for text between <text> tags
	for {
		startIdx := strings.Index(content, "<text")
		if startIdx == -1 {
			break
		}
		endTagStart := strings.Index(content[startIdx:], ">")
		if endTagStart == -1 {
			break
		}
		endIdx := strings.Index(content[startIdx:], "</text>")
		if endIdx == -1 {
			break
		}

		textContent := content[startIdx+endTagStart+1 : startIdx+endIdx]
		// Remove nested tags
		textContent = strings.ReplaceAll(textContent, "<tspan", "")
		textContent = strings.ReplaceAll(textContent, "</tspan>", "")
		for strings.Contains(textContent, ">") {
			start := strings.Index(textContent, ">")
			textContent = textContent[start+1:]
		}
		textContent = strings.TrimSpace(textContent)
		if textContent != "" {
			texts = append(texts, textContent)
		}

		content = content[startIdx+endIdx+7:]
	}

	return strings.Join(texts, " ")
}

func (c *VisionClient) getBasicImageInfo(path string, data []byte) string {
	ext := strings.ToLower(filepath.Ext(path))
	fileName := filepath.Base(path)

	info := fmt.Sprintf("Image: %s\nType: %s\nSize: %d bytes", fileName, ext, len(data))

	// Add base64 preview for small images
	if len(data) < 1000 {
		info += "\nBase64: " + base64.StdEncoding.EncodeToString(data)[:100] + "..."
	}

	return info
}
