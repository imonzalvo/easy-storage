package file

import (
	"errors"
)

// ErrFileNotFound is returned when a file cannot be found
var ErrFileNotFound = errors.New("file not found")

// Repository defines the interface for file data access
type Repository interface {
	Save(file *File) error
	FindByID(id string) (*File, error)
	FindByUserID(userID string, limit, offset int) ([]*File, error)
	Delete(id string) error
}
