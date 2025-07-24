package objestorage

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// GCS implements the ObjectStorage interface using Google Cloud Storage
type GCS struct {
	client     *storage.Client
	bucketName string
}

// NewGCS creates a new GCS instance
func NewGCS(bucketName string) (*GCS, error) {
	ctx := context.Background()

	// Create a new GCS client
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return &GCS{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// Upload stores a file in Google Cloud Storage
func (gcs *GCS) Upload(filename string, data io.Reader) (string, error) {
	ctx := context.Background()

	// Get a reference to the bucket and object
	obj := gcs.client.Bucket(gcs.bucketName).Object(filename)

	// Create a writer for the object
	writer := obj.NewWriter(ctx)
	defer writer.Close()

	// Copy data to GCS
	_, err := io.Copy(writer, data)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to GCS: %w", err)
	}

	// Return the filename as the identifier
	return filename, nil
}

// Retrieve gets a file from Google Cloud Storage
func (gcs *GCS) Retrieve(fileID string) ([]byte, error) {
	ctx := context.Background()

	// Get a reference to the object
	obj := gcs.client.Bucket(gcs.bucketName).Object(fileID)

	// Create a reader for the object
	reader, err := obj.NewReader(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil, fmt.Errorf("file not found: %s", fileID)
		}
		return nil, fmt.Errorf("failed to read file from GCS: %w", err)
	}
	defer reader.Close()

	// Read all data
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	return data, nil
}

// Close closes the GCS client
func (gcs *GCS) Close() error {
	return gcs.client.Close()
}
