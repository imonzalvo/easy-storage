package file

import (
	"time"
)

// File represents a file in the system
type File struct {
	ID          string
	Name        string
	Size        int64
	ContentType string
	Path        string
	UserID      string
	FolderID    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewFile creates a new file entity
func NewFile(name string, size int64, contentType, path, userID, folderID string) *File {
	now := time.Now()
	return &File{
		Name:        name,
		Size:        size,
		ContentType: contentType,
		Path:        path,
		UserID:      userID,
		FolderID:    folderID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
