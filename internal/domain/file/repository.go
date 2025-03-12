package file

import (
	"errors"
)

// ErrFileNotFound is returned when a file cannot be found
var ErrFileNotFound = errors.New("file not found")

// ErrInvalidFolder is returned when a folder doesn't exist or doesn't belong to the user
var ErrInvalidFolder = errors.New("invalid folder")

// Repository defines the interface for file data access
type Repository interface {
	Save(file *File) error
	FindByID(id string) (*File, error)
	FindByUserID(userID string, limit, offset int, sortBy, sortDir string) ([]*File, error)
	FindByUserIDAndFolder(userID string, folderID string) ([]*File, error)
	Delete(id string) error
	DeleteByFolder(folderID string) error
}
