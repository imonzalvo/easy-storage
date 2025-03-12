package share

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for share data access
type Repository interface {
	// Create stores a new share
	Create(ctx context.Context, share *Share) error

	// GetByID retrieves a share by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*Share, error)

	// GetByToken retrieves a share by its token
	GetByToken(ctx context.Context, token string) (*Share, error)

	// GetByResource retrieves all shares for a specific resource
	GetByResource(ctx context.Context, resourceID uuid.UUID, resourceType string) ([]*Share, error)

	// GetByOwner retrieves all shares created by a specific owner
	GetByOwner(ctx context.Context, ownerID uuid.UUID) ([]*Share, error)

	// GetByRecipient retrieves all shares shared with a specific recipient
	GetByRecipient(ctx context.Context, recipientID uuid.UUID) ([]*Share, error)

	// Update updates an existing share
	Update(ctx context.Context, share *Share) error

	// Delete removes a share
	Delete(ctx context.Context, id uuid.UUID) error
}
