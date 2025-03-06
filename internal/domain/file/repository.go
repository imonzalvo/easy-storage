package file

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for file persistence
type Repository interface {
	// Save persists a file entity
	Save(ctx context.Context, file *File) error

	// FindByID retrieves a file by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*File, error)

	// FindByUserAndFolder retrieves files for a user in a specific folder
	FindByUserAndFolder(ctx context.Context, userID uuid.UUID, folderID *uuid.UUID, limit, offset int) ([]*File, int64, error)

	// Delete removes a file
	Delete(ctx context.Context, id uuid.UUID) error

	// Update updates a file
	Update(ctx context.Context, file *File) error
	
	// GetTotalSizeByUser gets the total storage used by a user
	GetTotalSizeByUser(ctx context.Context, userID uuid.UUID) (int64, error)
}
