package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Folder represents a folder in the database
type Folder struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	Name      string `gorm:"not null"`
	ParentID  string `gorm:"type:uuid;default:null"`
	UserID    string `gorm:"type:uuid;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (f *Folder) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}
