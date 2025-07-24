package objestorage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Local implements the Storage interface using local filesystem
type Local struct {
	basePath string
}

// NewLocal creates a new Local instance
func NewLocal(basePath string) (*Local, error) {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &Local{
		basePath: basePath,
	}, nil
}

// Upload stores a file in the local filesystem
func (ls *Local) Upload(filename string, data io.Reader) (string, error) {
	// Create the full file path
	filePath := filepath.Join(ls.basePath, filename)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy data to file
	_, err = io.Copy(file, data)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return the filename as the identifier
	return filename, nil
}

// Retrieve gets a file from the local filesystem
func (ls *Local) Retrieve(fileID string) ([]byte, error) {
	// Create the full file path
	filePath := filepath.Join(ls.basePath, fileID)

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", fileID)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}
