package access

import (
	"context"
	"easy-storage/internal/domain/file"
	"easy-storage/internal/domain/share"
	"log"

	"github.com/google/uuid"
)

// Service handles file access control
type Service struct {
	fileService  *file.Service
	shareService *share.Service
}

// NewService creates a new file access service
func NewService(fileService *file.Service, shareService *share.Service) *Service {
	// Check if services are properly initialized
	if fileService == nil || shareService == nil {
		log.Printf("Error: fileService or shareService is nil")
	}
	return &Service{
		fileService:  fileService,
		shareService: shareService,
	}
}

// CheckFileAccess checks if a user has access to a file
func (s *Service) CheckFileAccess(ctx context.Context, fileID string, userID string) (bool, error) {
	// Parse UUIDs
	fileUUID, err := uuid.Parse(fileID)
	if err != nil {
		return false, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false, err
	}

	// Get file
	file, err := s.fileService.GetFile(fileID)
	if err != nil {
		return false, err
	}

	// Check if user is the owner
	if file.UserID == userID {
		return true, nil
	}

	// Check if file has been shared with user
	hasAccess, err := s.shareService.CheckAccessToResource(ctx, userUUID, fileUUID, "file")
	if err != nil {
		return false, err
	}

	return hasAccess, nil
}

// GetFileByShareToken gets a file using a share token
func (s *Service) GetFileByShareToken(ctx context.Context, token string, password string) (*file.File, error) {
	// Get share by token
	shareObj, err := s.shareService.GetResourceByToken(ctx, token, password)
	if err != nil {
		return nil, err
	}

	// Check if resource is a file
	if shareObj.ResourceType != "file" {
		return nil, ErrInvalidResourceType
	}

	// Get file
	fileID := shareObj.ResourceID.String()
	file, err := s.fileService.GetFile(fileID)
	if err != nil {
		return nil, err
	}

	return file, nil
}
