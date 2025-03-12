package user

import (
	"errors"
	"sync"
)

// ErrStorageQuotaExceeded is returned when a user tries to upload a file that would exceed their quota
var ErrStorageQuotaExceeded = errors.New("storage quota exceeded")

// StorageService provides user storage operations
type StorageService struct {
	repo Repository
	mu   sync.Mutex // Protects concurrent updates to storage statistics
}

// NewStorageService creates a new storage service
func NewStorageService(repo Repository) *StorageService {
	return &StorageService{
		repo: repo,
	}
}

// CheckQuota checks if a user has enough quota for a file of the given size
func (s *StorageService) CheckQuota(userID string, fileSize int64) (bool, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return false, err
	}

	return user.StorageUsed+fileSize <= user.StorageQuota, nil
}

// AddStorage increments a user's storage used
func (s *StorageService) AddStorage(userID string, size int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check quota first
	hasQuota, err := s.CheckQuota(userID, size)
	if err != nil {
		return err
	}
	if !hasQuota {
		return ErrStorageQuotaExceeded
	}

	return s.repo.IncrementStorageUsed(userID, size)
}

// RemoveStorage decrements a user's storage used
func (s *StorageService) RemoveStorage(userID string, size int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.DecrementStorageUsed(userID, size)
}

// GetStorageStats gets a user's storage statistics
func (s *StorageService) GetStorageStats(userID string) (quota int64, used int64, err error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return 0, 0, err
	}

	return user.StorageQuota, user.StorageUsed, nil
}

// RecalculateStorage recalculates a user's storage used based on their files
// This is useful for consistency checks or repairs
func (s *StorageService) RecalculateStorage(userID string, totalSize int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.UpdateStorageUsed(userID, totalSize)
}
