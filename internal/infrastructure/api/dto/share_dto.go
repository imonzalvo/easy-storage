package dto

import (
	"easy-storage/internal/domain/share"
	"time"

	"github.com/google/uuid"
)

// CreateShareRequest represents the request to create a new share
type CreateShareRequest struct {
	ResourceID   uuid.UUID  `json:"resource_id" validate:"required"`
	ResourceType string     `json:"resource_type" validate:"required,oneof=file folder"`
	ShareType    string     `json:"share_type" validate:"required,oneof=LINK USER"`
	Permission   string     `json:"permission" validate:"required,oneof=READ WRITE"`
	RecipientID  *uuid.UUID `json:"recipient_id,omitempty"`
	Password     *string    `json:"password,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
}

// ShareResponse represents the response for a share
type ShareResponse struct {
	ID           uuid.UUID  `json:"id"`
	ResourceID   uuid.UUID  `json:"resource_id"`
	ResourceType string     `json:"resource_type"`
	ShareType    string     `json:"share_type"`
	Permission   string     `json:"permission"`
	RecipientID  *uuid.UUID `json:"recipient_id,omitempty"`
	Token        string     `json:"token,omitempty"`
	HasPassword  bool       `json:"has_password"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	AccessCount  int        `json:"access_count"`
	LastAccessAt *time.Time `json:"last_access_at,omitempty"`
	IsRevoked    bool       `json:"is_revoked"`
	URL          string     `json:"url,omitempty"` // Full URL for link shares
}

// AccessShareRequest represents the request to access
// ... existing code ...

// AccessShareRequest represents the request to access a shared resource
type AccessShareRequest struct {
	Token    string  `json:"token" validate:"required"`
	Password *string `json:"password,omitempty"`
}

// RevokeShareRequest represents the request to revoke a share
type RevokeShareRequest struct {
	ShareID uuid.UUID `json:"share_id" validate:"required"`
}

// UpdateShareRequest represents the request to update a share
type UpdateShareRequest struct {
	ShareID   uuid.UUID  `json:"share_id" validate:"required"`
	Password  *string    `json:"password,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// ShareListResponse represents a list of shares
type ShareListResponse struct {
	Shares []ShareResponse `json:"shares"`
	Count  int             `json:"count"`
}

// MapDomainToResponse maps a domain share to a response DTO
func MapDomainToResponse(s *share.Share, baseURL string) ShareResponse {
	response := ShareResponse{
		ID:           s.ID,
		ResourceID:   s.ResourceID,
		ResourceType: s.ResourceType,
		ShareType:    string(s.Type),
		Permission:   string(s.Permission),
		RecipientID:  s.RecipientID,
		HasPassword:  s.Password != nil,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
		AccessCount:  s.AccessCount,
		LastAccessAt: s.LastAccessAt,
		IsRevoked:    s.IsRevoked,
	}

	// Only include token in response for link shares
	if s.Type == share.LinkShare {
		response.Token = s.Token
		// Construct the full URL for the share
		response.URL = baseURL + "/share/" + s.Token
	}

	return response
}

// MapDomainListToResponse maps a list of domain shares to a response DTO
func MapDomainListToResponse(shares []*share.Share, baseURL string) ShareListResponse {
	response := ShareListResponse{
		Shares: make([]ShareResponse, len(shares)),
		Count:  len(shares),
	}

	for i, s := range shares {
		response.Shares[i] = MapDomainToResponse(s, baseURL)
	}

	return response
}
