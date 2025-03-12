package share

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
)

// Service provides share-related operations
type Service struct {
	repo Repository
}

// NewService creates a new share service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateLinkShare creates a new link share for a resource
func (s *Service) CreateLinkShare(
	ctx context.Context,
	ownerID uuid.UUID,
	resourceID uuid.UUID,
	resourceType string,
	permission SharePermission,
) (*Share, error) {
	share := NewShare(ownerID, resourceID, resourceType, LinkShare, permission)

	// Generate a secure random token
	token, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	share.SetToken(token)

	if err := s.repo.Create(ctx, share); err != nil {
		return nil, err
	}

	return share, nil
}

// CreateUserShare creates a new direct user share
func (s *Service) CreateUserShare(
	ctx context.Context,
	ownerID uuid.UUID,
	resourceID uuid.UUID,
	resourceType string,
	recipientID uuid.UUID,
	permission SharePermission,
) (*Share, error) {
	share := NewShare(ownerID, resourceID, resourceType, UserShare, permission)
	share.SetRecipient(recipientID)

	if err := s.repo.Create(ctx, share); err != nil {
		return nil, err
	}

	return share, nil
}

// GetShareByID retrieves a share by its ID
func (s *Service) GetShareByID(ctx context.Context, id uuid.UUID) (*Share, error) {
	return s.repo.GetByID(ctx, id)
}

// GetShareByToken retrieves a share by its token
func (s *Service) GetShareByToken(ctx context.Context, token string) (*Share, error) {
	return s.repo.GetByToken(ctx, token)
}

// SetShareExpiration sets an expiration date for a share
func (s *Service) SetShareExpiration(ctx context.Context, shareID uuid.UUID, expiresAt time.Time) error {
	share, err := s.repo.GetByID(ctx, shareID)
	if err != nil {
		return err
	}

	share.SetExpiration(expiresAt)
	return s.repo.Update(ctx, share)
}

// SetSharePassword sets a password for a share
func (s *Service) SetSharePassword(ctx context.Context, shareID uuid.UUID, password string) error {
	share, err := s.repo.GetByID(ctx, shareID)
	if err != nil {
		return err
	}

	share.SetPassword(password)
	return s.repo.Update(ctx, share)
}

// RevokeShare revokes access to a share
func (s *Service) RevokeShare(ctx context.Context, shareID uuid.UUID) error {
	share, err := s.repo.GetByID(ctx, shareID)
	if err != nil {
		return err
	}

	share.Revoke()
	return s.repo.Update(ctx, share)
}

// RecordShareAccess records an access to a share
func (s *Service) RecordShareAccess(ctx context.Context, shareID uuid.UUID) error {
	share, err := s.repo.GetByID(ctx, shareID)
	if err != nil {
		return err
	}

	share.RecordAccess()
	return s.repo.Update(ctx, share)
}

// ListSharesByOwner lists all shares created by an owner
func (s *Service) ListSharesByOwner(ctx context.Context, ownerID uuid.UUID) ([]*Share, error) {
	return s.repo.GetByOwner(ctx, ownerID)
}

// ListSharesByResource lists all shares for a specific resource
func (s *Service) ListSharesByResource(ctx context.Context, resourceID uuid.UUID, resourceType string) ([]*Share, error) {
	return s.repo.GetByResource(ctx, resourceID, resourceType)
}

// ListSharesWithUser lists all shares shared with a specific user
func (s *Service) ListSharesWithUser(ctx context.Context, userID uuid.UUID) ([]*Share, error) {
	return s.repo.GetByRecipient(ctx, userID)
}

// DeleteShare permanently deletes a share
func (s *Service) DeleteShare(ctx context.Context, shareID uuid.UUID) error {
	return s.repo.Delete(ctx, shareID)
}

// ValidateSharePassword checks if the provided password is valid for a share
func (s *Service) ValidateSharePassword(ctx context.Context, shareID uuid.UUID, password string) (bool, error) {
	share, err := s.repo.GetByID(ctx, shareID)
	if err != nil {
		return false, err
	}

	// If no password is set, any password is invalid
	if share.Password == nil {
		return false, nil
	}

	return *share.Password == password, nil
}

// GetShareURL generates a full URL for a share based on the base URL
func (s *Service) GetShareURL(share *Share, baseURL string) (string, error) {
	if share.Type != LinkShare || share.Token == "" {
		return "", ErrInvalidShareType
	}

	return baseURL + "/share/" + share.Token, nil
}

// CheckAccessToResource checks if a user has access to a resource through shares
func (s *Service) CheckAccessToResource(
	ctx context.Context,
	userID uuid.UUID,
	resourceID uuid.UUID,
	resourceType string,
) (bool, error) {
	// Get all shares for this resource
	shares, err := s.repo.GetByResource(ctx, resourceID, resourceType)
	if err != nil {
		return false, err
	}

	// Check if any share gives this user access
	for _, share := range shares {
		// Skip if share is revoked or expired
		if share.IsRevoked || (share.ExpiresAt != nil && share.IsExpired()) {
			continue
		}

		// Check if this is a user share directly with this user
		if share.Type == UserShare && share.RecipientID != nil && *share.RecipientID == userID {
			return true, nil
		}
	}

	return false, nil
}

// GetResourceByToken retrieves resource information from a share token
func (s *Service) GetResourceByToken(
	ctx context.Context,
	token string,
	password string,
) (*Share, error) {
	// Get share by token
	share, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Check if share is accessible
	if !share.IsAccessible() {
		if share.IsRevoked {
			return nil, ErrShareRevoked
		}
		if share.IsExpired() {
			return nil, ErrShareExpired
		}
	}

	// Check password if required
	if share.Password != nil && *share.Password != password {
		return nil, ErrInvalidPassword
	}

	// Record access
	share.RecordAccess()
	if err := s.repo.Update(ctx, share); err != nil {
		// Log error but continue
		// logger.Warn("Failed to record share access", "error", err)
	}

	return share, nil
}

// Helper function to generate a secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
