package file

import (
	"time"

	"github.com/google/uuid"
)

// File represents a file entity in the domain
type File struct {
	ID           uuid.UUID
	Name         string
	Size         int64
	ContentType  string
	Path         string
	StorageKey   string
	FolderID     *uuid.UUID
	UserID       uuid.UUID
	IsPublic     bool
	ThumbnailKey string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewFile creates a new file entity
func NewFile(name string, size int64, contentType string, userID uuid.UUID, folderID *uuid.UUID) *File {
	return &File{
		ID:          uuid.New(),
		Name:        name,
		Size:        size,
		ContentType: contentType,
		UserID:      userID,
		FolderID:    folderID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsPublic:    false,
	}
}

// SetStorageKey sets the storage key for the file
func (f *File) SetStorageKey(key string) {
	f.StorageKey = key
	f.UpdatedAt = time.Now()
}

// SetPublic marks the file as public or private
func (f *File) SetPublic(isPublic bool) {
	f.IsPublic = isPublic
	f.UpdatedAt = time.Now()
}

// Rename changes the file name
func (f *File) Rename(newName string) {
	f.Name = newName
	f.UpdatedAt = time.Now()
}

// Move changes the folder ID for the file
func (f *File) Move(folderID *uuid.UUID) {
	f.FolderID = folderID
	f.UpdatedAt = time.Now()
}
