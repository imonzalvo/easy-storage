package share

import (
	"time"

	"github.com/google/uuid"
)

// ShareType defines the type of sharing
type ShareType string

const (
	// LinkShare represents sharing via a public link
	LinkShare ShareType = "LINK"
	// UserShare represents sharing directly with another user
	UserShare ShareType = "USER"
)

// SharePermission defines the permission level for a share
type SharePermission string

const (
	// ReadOnly permission allows only viewing files
	ReadOnly SharePermission = "READ"
	// ReadWrite permission allows viewing and modifying files
	ReadWrite SharePermission = "WRITE"
)

// Share represents a sharing entity in the system
type Share struct {
	ID           uuid.UUID       `json:"id"`
	OwnerID      uuid.UUID       `json:"owner_id"`
	ResourceID   uuid.UUID       `json:"resource_id"`   // Can be file or folder ID
	ResourceType string          `json:"resource_type"` // "file" or "folder"
	Type         ShareType       `json:"type"`
	Permission   SharePermission `json:"permission"`
	RecipientID  *uuid.UUID      `json:"recipient_id,omitempty"` // Only for UserShare
	Token        string          `json:"token,omitempty"`        // For LinkShare
	Password     *string         `json:"password,omitempty"`     // Optional password protection
	ExpiresAt    *time.Time      `json:"expires_at,omitempty"`   // Optional expiration
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	AccessCount  int             `json:"access_count"` // Track number of accesses
	LastAccessAt *time.Time      `json:"last_access_at,omitempty"`
	IsRevoked    bool            `json:"is_revoked"`
}

// NewShare creates a new share entity
func NewShare(
	ownerID uuid.UUID,
	resourceID uuid.UUID,
	resourceType string,
	shareType ShareType,
	permission SharePermission,
) *Share {
	return &Share{
		ID:           uuid.New(),
		OwnerID:      ownerID,
		ResourceID:   resourceID,
		ResourceType: resourceType,
		Type:         shareType,
		Permission:   permission,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		AccessCount:  0,
		IsRevoked:    false,
	}
}

// SetPassword adds password protection to the share
func (s *Share) SetPassword(password string) {
	s.Password = &password
	s.UpdatedAt = time.Now()
}

// SetExpiration sets an expiration date for the share
func (s *Share) SetExpiration(expiresAt time.Time) {
	s.ExpiresAt = &expiresAt
	s.UpdatedAt = time.Now()
}

// SetRecipient sets a specific recipient for user sharing
func (s *Share) SetRecipient(recipientID uuid.UUID) {
	s.RecipientID = &recipientID
	s.UpdatedAt = time.Now()
}

// SetToken sets the access token for link sharing
func (s *Share) SetToken(token string) {
	s.Token = token
	s.UpdatedAt = time.Now()
}

// Revoke revokes access to this share
func (s *Share) Revoke() {
	s.IsRevoked = true
	s.UpdatedAt = time.Now()
}

// RecordAccess records an access to this share
func (s *Share) RecordAccess() {
	s.AccessCount++
	now := time.Now()
	s.LastAccessAt = &now
}

// IsExpired checks if the share has expired
func (s *Share) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*s.ExpiresAt)
}

// IsAccessible checks if the share is currently accessible
func (s *Share) IsAccessible() bool {
	return !s.IsRevoked && !s.IsExpired()
}
