package file

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"easy-storage/internal/domain/user"
)

// StorageProvider defines the interface for file storage operations
type StorageProvider interface {
	// Upload stores a file
	Upload(ctx context.Context, key string, file multipart.File, size int64, contentType string) error

	// Download retrieves a file
	Download(ctx context.Context, key string) ([]byte, error)

	// Delete removes a file
	Delete(ctx context.Context, key string) error
}

// Service handles file domain logic
type Service struct {
	repo           Repository
	storageProvider StorageProvider
	userRepo        user.Repository
}

// NewService creates a new file service
func NewService(repo Repository, storageProvider StorageProvider, userRepo user.Repository) *Service {
	return &Service{
		repo:           repo,
		storageProvider: storageProvider,
		userRepo:        userRepo,
	}
}

// UploadFile uploads a new file
func (s *Service) UploadFile(ctx context.Context, name string, fileHeader *multipart.FileHeader, 
	contentType string, userID uuid.UUID, folderID *uuid.UUID) (*File, error) {
	
	// Check user storage quota
	usr, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	
	totalUsed, err := s.repo.GetTotalSizeByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user storage usage: %w", err)
	}
	
	if usr.StorageQuota > 0 && totalUsed+fileHeader.Size > usr.StorageQuota {
		return nil, ErrStorageQuotaExceeded
	}
	
	// Create file entity
	file := NewFile(name, fileHeader.Size, contentType, userID, folderID)
	
	// Generate storage key
	storageKey := fmt.Sprintf("%s/%s/%s", userID.String(), time.Now().Format("2006/01/02"), file.ID.String())
	file.SetStorageKey(storageKey)
	
	// Open file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()
	
	// Upload to storage provider
	if err := s.storageProvider.Upload(ctx, storageKey, src, fileHeader.Size, contentType); err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}
	
	// Save to repository
	if err := s.repo.Save(ctx, file); err != nil {
		// Attempt to delete from storage on failure
		_ = s.storageProvider.Delete(ctx, storageKey)
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}
	
	return file, nil
}

// GetFile retrieves a file by ID
func (s *Service) GetFile(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*File, error) {
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Check if user has access to this file
	if file.UserID != userID && !file.IsPublic {
		return nil, ErrAccessDenied
	}
	
	return file, nil
}

// DeleteFile removes a file
func (s *Service) DeleteFile(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// Check if user has access to this file
	if file.UserID != userID {
		return ErrAccessDenied
	}
	
	// Delete from storage
	if err := s.storageProvider.Delete(ctx, file.StorageKey); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}
	
	// Delete from repository
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete file metadata: %w", err)
	}
	
	return nil
}