// internal/infrastructure/persistence/gorm/models/user.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the user model in the database
type User struct {
	ID           string `gorm:"primaryKey;type:uuid"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Name         string `gorm:"not null"`
	StorageQuota int64  `gorm:"default:5368709120"` // Default 5GB in bytes
	StorageUsed  int64  `gorm:"default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
