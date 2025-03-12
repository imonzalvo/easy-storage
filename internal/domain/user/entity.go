package user

import "time"

// User represents the user entity
type User struct {
	ID           string
	Email        string
	Password     string
	Name         string
	StorageQuota int64 // Default 5GB in bytes
	StorageUsed  int64 // Current storage used in bytes
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser creates a new user entity
func NewUser(email, password, name string) *User {
	now := time.Now()
	return &User{
		Email:        email,
		Password:     password,
		Name:         name,
		StorageQuota: 5 * 1024 * 1024 * 1024, // 5GB default
		StorageUsed:  0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
