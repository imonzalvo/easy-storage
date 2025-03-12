package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Share represents the database model for shares
type Share struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key"`
	OwnerID      uuid.UUID  `gorm:"type:uuid;index"`
	ResourceID   uuid.UUID  `gorm:"type:uuid;index"`
	ResourceType string     `gorm:"type:varchar(10);index"`
	Type         string     `gorm:"type:varchar(10)"`
	Permission   string     `gorm:"type:varchar(10)"`
	RecipientID  *uuid.UUID `gorm:"type:uuid;index;null"`
	Token        string     `gorm:"type:varchar(255);index;null"`
	Password     *string    `gorm:"type:varchar(255);null"`
	ExpiresAt    *time.Time `gorm:"null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	AccessCount  int        `gorm:"default:0"`
	LastAccessAt *time.Time `gorm:"null"`
	IsRevoked    bool       `gorm:"default:false"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (s *Share) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
