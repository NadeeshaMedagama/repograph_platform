package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

// ComputeFileHash computes SHA256 hash of a file
func ComputeFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to compute hash: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetMimeType returns the MIME type of a file
func GetMimeType(filePath string) string {
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream"
	}
	return mimeType
}

// GetFileExtension returns the file extension without the dot
func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext != "" {
		return strings.ToLower(strings.TrimPrefix(ext, "."))
	}
	return ""
}

// IsImageFile checks if a file is an image based on extension
func IsImageFile(filename string) bool {
	imageExtensions := []string{"png", "jpg", "jpeg", "gif", "bmp", "svg", "webp"}
	ext := GetFileExtension(filename)
	for _, imgExt := range imageExtensions {
		if ext == imgExt {
			return true
		}
	}
	return false
}

// IsDiagramFile checks if a file is a diagram based on extension
func IsDiagramFile(filename string) bool {
	diagramExtensions := []string{"drawio", "excalidraw"}
	ext := GetFileExtension(filename)
	for _, diagExt := range diagramExtensions {
		if ext == diagExt {
			return true
		}
	}
	return false
}

// IsDocumentFile checks if a file is a document based on extension
func IsDocumentFile(filename string) bool {
	docExtensions := []string{"docx", "pdf", "pptx", "odt", "txt", "md"}
	ext := GetFileExtension(filename)
	for _, docExt := range docExtensions {
		if ext == docExt {
			return true
		}
	}
	return false
}

// IsSpreadsheetFile checks if a file is a spreadsheet based on extension
func IsSpreadsheetFile(filename string) bool {
	sheetExtensions := []string{"xlsx", "xls", "csv"}
	ext := GetFileExtension(filename)
	for _, sheetExt := range sheetExtensions {
		if ext == sheetExt {
			return true
		}
	}
	return false
}

// IsCodeFile checks if a file is a code file based on extension
func IsCodeFile(filename string) bool {
	codeExtensions := []string{
		"go", "py", "js", "ts", "java", "c", "cpp", "h", "hpp",
		"rs", "rb", "php", "swift", "kt", "scala", "r", "sql",
		"sh", "bash", "ps1", "dart", "lua", "perl", "groovy",
	}
	ext := GetFileExtension(filename)
	for _, codeExt := range codeExtensions {
		if ext == codeExt {
			return true
		}
	}
	return false
}

// IsStructuredFile checks if a file is a structured data file based on extension
func IsStructuredFile(filename string) bool {
	structuredExtensions := []string{"json", "yaml", "yml", "xml", "toml", "graphql"}
	ext := GetFileExtension(filename)
	for _, structExt := range structuredExtensions {
		if ext == structExt {
			return true
		}
	}
	return false
}

// GetFileCategory returns the category of a file based on its extension
func GetFileCategory(filename string) string {
	switch {
	case IsImageFile(filename):
		return "image"
	case IsDiagramFile(filename):
		return "diagram"
	case IsDocumentFile(filename):
		return "document"
	case IsSpreadsheetFile(filename):
		return "spreadsheet"
	case IsCodeFile(filename):
		return "code"
	case IsStructuredFile(filename):
		return "structured"
	default:
		return "unknown"
	}
}

// FileExists checks if a file exists
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists
func DirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(dirPath string) error {
	if !DirExists(dirPath) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// ReadFile reads the entire contents of a file
func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// WriteFile writes data to a file
func WriteFile(filePath string, data []byte) error {
	return os.WriteFile(filePath, data, 0644)
}

// SanitizeFileName removes or replaces invalid characters in a filename
func SanitizeFileName(name string) string {
	// Replace invalid characters with underscores
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	sanitized := name
	for _, char := range invalidChars {
		sanitized = strings.ReplaceAll(sanitized, char, "_")
	}
	return sanitized
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
