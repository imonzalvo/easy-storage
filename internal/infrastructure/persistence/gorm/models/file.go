package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// File represents a file in the database
type File struct {
	ID          string `gorm:"primaryKey;type:uuid"`
	Name        string `gorm:"not null"`
	Size        int64  `gorm:"not null"`
	ContentType string `gorm:"not null"`
	Path        string `gorm:"not null"` // Path in the storage system
	UserID      string `gorm:"type:uuid;not null"`
	FolderID    string `gorm:"type:uuid;default:null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}
