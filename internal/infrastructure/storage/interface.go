// internal/infrastructure/storage/interface.go
package storage

import (
	"io"
)

// Provider defines the interface for file storage
type Provider interface {
	// Upload uploads a file to storage and returns its path
	Upload(filename string, contentType string, file io.Reader) (string, error)

	// Download downloads a file from storage
	Download(path string) (io.ReadCloser, error)

	// Delete deletes a file from storage
	Delete(path string) error

	// GetSignedURL generates a presigned URL for downloading a file
	// expiryTime is the duration in seconds for which the URL will be valid
	GetSignedURL(path string, expiryTime int64) (string, error)
}
