package file

import (
	"easy-storage/internal/domain/common"
	"easy-storage/internal/domain/user"
	"io"
	"log"
)

// StorageProvider defines the interface for file storage operations
type StorageProvider interface {
	Upload(filename string, contentType string, file io.Reader) (string, error)
	Download(path string) (io.ReadCloser, error)
	Delete(path string) error
	GetSignedURL(path string, expiryTime int64) (string, error)
}

// Service provides file operations
type Service struct {
	repo            Repository
	folderValidator common.FolderValidator
	storage         StorageProvider
	userStorage     *user.StorageService
}

// NewService creates a new file service
func NewService(repo Repository, folderValidator common.FolderValidator, storage StorageProvider, userStorage *user.StorageService) *Service {
	return &Service{
		repo:            repo,
		folderValidator: folderValidator,
		storage:         storage,
		userStorage:     userStorage,
	}
}

// UploadFile uploads a file to storage and saves metadata
func (s *Service) UploadFile(filename string, size int64, contentType string, fileContent io.Reader, userID, folderID string) (*File, error) {
	// Validate folder ownership if folderID is provided
	if folderID != "" {
		// Check if folder exists and belongs to the user
		belongs, err := s.folderValidator.BelongsToUser(folderID, userID)
		if err != nil {
			return nil, err
		}
		if !belongs {
			return nil, ErrInvalidFolder
		}
	}

	// Check if user has enough storage quota
	if s.userStorage != nil {
		if err := s.userStorage.AddStorage(userID, size); err != nil {
			if err == user.ErrStorageQuotaExceeded {
				return nil, err
			}
			return nil, err
		}
	}

	// Upload file to storage
	path, err := s.storage.Upload(filename, contentType, fileContent)
	if err != nil {
		// Rollback storage usage increment if upload fails
		if s.userStorage != nil {
			_ = s.userStorage.RemoveStorage(userID, size)
		}
		return nil, err
	}

	// Create file entity
	file := NewFile(filename, size, contentType, path, userID, folderID)

	// Save file metadata to repository
	if err := s.repo.Save(file); err != nil {
		// Try to clean up the stored file if metadata save fails
		_ = s.storage.Delete(path)
		// Rollback storage usage increment if save fails
		if s.userStorage != nil {
			_ = s.userStorage.RemoveStorage(userID, size)
		}
		return nil, err
	}

	return file, nil
}

// GetFile retrieves a file by ID
func (s *Service) GetFile(id string) (*File, error) {
	return s.repo.FindByID(id)
}

// GetFileContent gets the content of a file
func (s *Service) GetFileContent(file *File) (io.ReadCloser, error) {
	return s.storage.Download(file.Path)
}

// DeleteFile deletes a file
func (s *Service) DeleteFile(id string) error {
	// Get file first to get the path
	file, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Delete from repository
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Update user storage statistics
	if s.userStorage != nil {
		if err := s.userStorage.RemoveStorage(file.UserID, file.Size); err != nil {
			log.Printf("Error updating storage statistics: %v", err)
			// Continue with deletion even if stats update fails
		}
	}

	// Delete from storage
	return s.storage.Delete(file.Path)
}

// ListUserFiles lists files for a user
func (s *Service) ListUserFiles(userID string, limit, offset int, sortBy, sortDir string) ([]*File, error) {
	return s.repo.FindByUserID(userID, limit, offset, sortBy, sortDir)
}

// GetFileSignedURL returns a signed URL for a file
func (s *Service) GetFileSignedURL(file *File, expiryTime int64) (string, error) {
	return s.storage.GetSignedURL(file.Path, expiryTime)
}

// ListFilesInFolder lists files for a user in a specific folder
func (s *Service) ListFilesInFolder(userID string, folderID string) ([]*File, error) {
	// If folderID is empty, list files in the root folder
	if folderID == "" {
		return s.repo.FindByUserIDAndFolder(userID, "")
	}

	// Validate folder ownership
	belongs, err := s.folderValidator.BelongsToUser(folderID, userID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, ErrInvalidFolder
	}

	return s.repo.FindByUserIDAndFolder(userID, folderID)
}

// DeleteByFolder deletes all files in a folder
func (s *Service) DeleteByFolder(userID, folderID string) error {
	// Get all files in the folder
	files, err := s.repo.FindByUserIDAndFolder(userID, folderID)
	if err != nil {
		return err
	}

	// Delete each file individually to ensure proper storage cleanup
	for _, file := range files {
		// Delete from storage
		if err := s.storage.Delete(file.Path); err != nil {
			// Log error but continue with other deletions
			// We don't want to stop the process if one file fails to delete
			// from storage, but we should log it for investigation
			log.Printf("Error deleting file from storage: %v", err)
		}

		// Delete from repository
		if err := s.repo.Delete(file.ID); err != nil {
			return err
		}

		// Update user storage statistics
		if s.userStorage != nil {
			if err := s.userStorage.RemoveStorage(userID, file.Size); err != nil {
				log.Printf("Error updating storage statistics: %v", err)
				// Continue with deletion even if stats update fails
			}
		}
	}

	return nil
}
