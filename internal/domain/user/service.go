// internal/domain/user/service.go
package user

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ErrUserNotFound is returned when a user cannot be found
var ErrUserNotFound = errors.New("user not found")

// ErrInvalidCredentials is returned when login credentials are invalid
var ErrInvalidCredentials = errors.New("invalid credentials")

// Service provides user operations
type Service struct {
	repo           Repository
	storageService *StorageService
}

// NewService creates a new user service
func NewService(repo Repository) *Service {
	storageService := NewStorageService(repo)
	return &Service{
		repo:           repo,
		storageService: storageService,
	}
}

// RegisterUser registers a new user
func (s *Service) RegisterUser(email, password, name string) (*User, error) {
	// Check if user with email already exists
	existingUser, err := s.repo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already in use")
	} else if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        email,
		Password:     string(hashedPassword),
		Name:         name,
		StorageQuota: 5 * 1024 * 1024 * 1024, // 5GB default
		StorageUsed:  0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Authenticate authenticates a user by email and password
func (s *Service) Authenticate(email, password string) (*User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(id string) (*User, error) {
	return s.repo.FindByID(id)
}

// UpdateUser updates a user's information
func (s *Service) UpdateUser(user *User) error {
	return s.repo.Update(user)
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(id string) error {
	return s.repo.Delete(id)
}

// ChangePassword changes a user's password
func (s *Service) ChangePassword(userID, currentPassword, newPassword string) error {
	// Get the user
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update user with new password
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.repo.Update(user)
}

// GetStorageService returns the storage service
func (s *Service) GetStorageService() *StorageService {
	return s.storageService
}

// GetStorageStats gets a user's storage statistics
func (s *Service) GetStorageStats(userID string) (quota int64, used int64, err error) {
	return s.storageService.GetStorageStats(userID)
}
