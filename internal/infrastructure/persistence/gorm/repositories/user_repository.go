// internal/infrastructure/persistence/gorm/repositories/user_repository.go
package repositories

import (
	"errors"

	"easy-storage/internal/domain/user"
	"easy-storage/internal/infrastructure/persistence/gorm/models"

	"gorm.io/gorm"
)

// GormUserRepository implements the user.Repository interface using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new user repository
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Save creates or updates a user in the database
func (r *GormUserRepository) Save(u *user.User) error {
	userModel := &models.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.Password,
		Name:         u.Name,
		StorageQuota: u.StorageQuota,
		StorageUsed:  u.StorageUsed,
	}

	if err := r.db.Save(userModel).Error; err != nil {
		return err
	}

	u.ID = userModel.ID
	return nil
}

// FindByID finds a user by ID
func (r *GormUserRepository) FindByID(id string) (*user.User, error) {
	var userModel models.User
	if err := r.db.First(&userModel, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return &user.User{
		ID:           userModel.ID,
		Email:        userModel.Email,
		Password:     userModel.PasswordHash,
		Name:         userModel.Name,
		StorageQuota: userModel.StorageQuota,
		StorageUsed:  userModel.StorageUsed,
		CreatedAt:    userModel.CreatedAt,
		UpdatedAt:    userModel.UpdatedAt,
	}, nil
}

// FindByEmail finds a user by email
func (r *GormUserRepository) FindByEmail(email string) (*user.User, error) {
	var userModel models.User
	if err := r.db.First(&userModel, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return &user.User{
		ID:           userModel.ID,
		Email:        userModel.Email,
		Password:     userModel.PasswordHash,
		Name:         userModel.Name,
		StorageQuota: userModel.StorageQuota,
		StorageUsed:  userModel.StorageUsed,
		CreatedAt:    userModel.CreatedAt,
		UpdatedAt:    userModel.UpdatedAt,
	}, nil
}

// Update updates a user
func (r *GormUserRepository) Update(u *user.User) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", u.ID).
		Updates(map[string]interface{}{
			"email":         u.Email,
			"password_hash": u.Password,
			"name":          u.Name,
			"storage_quota": u.StorageQuota,
			"storage_used":  u.StorageUsed,
		}).Error
}

// Delete deletes a user
func (r *GormUserRepository) Delete(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

// UpdateStorageUsed updates only the storage used field for a user
func (r *GormUserRepository) UpdateStorageUsed(userID string, storageUsed int64) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("storage_used", storageUsed).Error
}

// IncrementStorageUsed increments the storage used field for a user
func (r *GormUserRepository) IncrementStorageUsed(userID string, size int64) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("storage_used", gorm.Expr("storage_used + ?", size)).Error
}

// DecrementStorageUsed decrements the storage used field for a user
func (r *GormUserRepository) DecrementStorageUsed(userID string, size int64) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("storage_used", gorm.Expr("GREATEST(storage_used - ?, 0)", size)).Error
}
