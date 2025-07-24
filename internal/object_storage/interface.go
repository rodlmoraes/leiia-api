package objestorage

import "io"

// ObjectStorage defines the interface for file storage operations
type ObjectStorage interface {
	// Upload stores a file and returns a unique identifier for it
	Upload(filename string, data io.Reader) (string, error)

	// Retrieve gets a file by its identifier and returns the data
	Retrieve(fileID string) ([]byte, error)
}
