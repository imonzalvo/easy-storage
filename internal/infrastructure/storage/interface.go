package storage

import (
	"context"
	"io"
)

// Storage defines the interface for storage operations
type Storage interface {
	// PutObject uploads an object to storage
	PutObject(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	
	// GetObject retrieves an object from storage
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	
	// DeleteObject removes an object from storage
	DeleteObject(ctx context.Context, key string) error
	
	// GetObjectURL gets a presigned URL for direct access
	GetObjectURL(ctx context.Context, key string, expires int64) (string, error)
}
