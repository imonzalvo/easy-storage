package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"easy-storage/internal/domain/share"
	"easy-storage/internal/infrastructure/persistence/gorm/models"
)

// ShareRepository implements the share.Repository interface using GORM
type ShareRepository struct {
	db *gorm.DB
}

// NewShareRepository creates a new share repository
func NewShareRepository(db *gorm.DB) *ShareRepository {
	return &ShareRepository{
		db: db,
	}
}

// Create stores a new share
func (r *ShareRepository) Create(ctx context.Context, s *share.Share) error {
	model := mapDomainToModel(s)
	result := r.db.WithContext(ctx).Create(model)
	return result.Error
}

// GetByID retrieves a share by its ID
func (r *ShareRepository) GetByID(ctx context.Context, id uuid.UUID) (*share.Share, error) {
	var model models.Share
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, share.ErrShareNotFound
		}
		return nil, result.Error
	}
	return mapModelToDomain(&model), nil
}

// GetByToken retrieves a share by its token
func (r *ShareRepository) GetByToken(ctx context.Context, token string) (*share.Share, error) {
	var model models.Share
	result := r.db.WithContext(ctx).Where("token = ?", token).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, share.ErrShareNotFound
		}
		return nil, result.Error
	}
	return mapModelToDomain(&model), nil
}

// GetByResource retrieves all shares for a specific resource
func (r *ShareRepository) GetByResource(ctx context.Context, resourceID uuid.UUID, resourceType string) ([]*share.Share, error) {
	var models []models.Share
	result := r.db.WithContext(ctx).Where("resource_id = ? AND resource_type = ?", resourceID, resourceType).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}
	return mapModelsToDomain(models), nil
}

// GetByOwner retrieves all shares created by
// ... existing code ...

// GetByOwner retrieves all shares created by a specific owner
func (r *ShareRepository) GetByOwner(ctx context.Context, ownerID uuid.UUID) ([]*share.Share, error) {
	var models []models.Share
	result := r.db.WithContext(ctx).Where("owner_id = ?", ownerID).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}
	return mapModelsToDomain(models), nil
}

// GetByRecipient retrieves all shares shared with a specific recipient
func (r *ShareRepository) GetByRecipient(ctx context.Context, recipientID uuid.UUID) ([]*share.Share, error) {
	var models []models.Share
	result := r.db.WithContext(ctx).Where("recipient_id = ?", recipientID).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}
	return mapModelsToDomain(models), nil
}

// Update updates an existing share
func (r *ShareRepository) Update(ctx context.Context, s *share.Share) error {
	model := mapDomainToModel(s)
	result := r.db.WithContext(ctx).Save(model)
	return result.Error
}

// Delete removes a share
func (r *ShareRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Share{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return share.ErrShareNotFound
	}
	return nil
}

// Helper functions to map between domain and model
func mapDomainToModel(s *share.Share) *models.Share {
	return &models.Share{
		ID:           s.ID,
		OwnerID:      s.OwnerID,
		ResourceID:   s.ResourceID,
		ResourceType: s.ResourceType,
		Type:         string(s.Type),
		Permission:   string(s.Permission),
		RecipientID:  s.RecipientID,
		Token:        s.Token,
		Password:     s.Password,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
		AccessCount:  s.AccessCount,
		LastAccessAt: s.LastAccessAt,
		IsRevoked:    s.IsRevoked,
	}
}

func mapModelToDomain(m *models.Share) *share.Share {
	return &share.Share{
		ID:           m.ID,
		OwnerID:      m.OwnerID,
		ResourceID:   m.ResourceID,
		ResourceType: m.ResourceType,
		Type:         share.ShareType(m.Type),
		Permission:   share.SharePermission(m.Permission),
		RecipientID:  m.RecipientID,
		Token:        m.Token,
		Password:     m.Password,
		ExpiresAt:    m.ExpiresAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		AccessCount:  m.AccessCount,
		LastAccessAt: m.LastAccessAt,
		IsRevoked:    m.IsRevoked,
	}
}

func mapModelsToDomain(models []models.Share) []*share.Share {
	shares := make([]*share.Share, len(models))
	for i, model := range models {
		shares[i] = mapModelToDomain(&model)
	}
	return shares
}
