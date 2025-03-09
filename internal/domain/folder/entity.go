package folder

import (
	"time"
)

// Folder represents a folder in the system
type Folder struct {
	ID        string
	Name      string
	ParentID  string // ID of parent folder, empty if root folder
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewFolder creates a new folder entity
func NewFolder(name string, parentID string, userID string) *Folder {
	now := time.Now()
	return &Folder{
		Name:      name,
		ParentID:  parentID,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
