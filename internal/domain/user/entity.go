package user

import "time"

// User represents the user entity
type User struct {
	ID        string
	Email     string
	Password  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user entity
func NewUser(email, password, name string) *User {
	now := time.Now()
	return &User{
		Email:     email,
		Password:  password,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
