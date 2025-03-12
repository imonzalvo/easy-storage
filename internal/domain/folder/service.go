package folder

import "easy-storage/internal/domain/file"

// Service provides folder operations
type Service struct {
	repo        Repository
	fileService *file.Service
}

// NewService creates a new folder service
func NewService(repo Repository, fileService *file.Service) *Service {
	return &Service{
		repo:        repo,
		fileService: fileService,
	}
}

// CreateFolder creates a new folder
func (s *Service) CreateFolder(name string, parentID string, userID string) (*Folder, error) {
	// If parentID is provided, verify it exists and belongs to user
	if parentID != "" {
		belongs, err := s.repo.BelongsToUser(parentID, userID)
		if err != nil {
			return nil, err
		}
		if !belongs {
			return nil, ErrInvalidParent
		}
	}

	folder := NewFolder(name, parentID, userID)
	if err := s.repo.Save(folder); err != nil {
		return nil, err
	}

	return folder, nil
}

// GetFolder retrieves a folder by ID
func (s *Service) GetFolder(id string) (*Folder, error) {
	return s.repo.FindByID(id)
}

// ListUserFolders lists all folders for a user
func (s *Service) ListUserFolders(userID string) ([]*Folder, error) {
	return s.repo.FindByUserID(userID)
}

// // DeleteFolder deletes a folder
// func (s *Service) DeleteFolder(id string) error {
// 	return s.repo.Delete(id)
// }

// ListFoldersByParent returns all folders for a user within a specific parent folder
// If parentID is empty, it returns root-level folders
func (s *Service) ListFoldersByParent(userID string, parentID string) ([]Folder, error) {
	return s.repo.FindByUserAndParent(userID, parentID)
}

// ListFoldersByParentPaginated returns paginated folders for a user within a specific parent folder
// If parentID is empty, it returns root-level folders
// It also returns the total count of folders and calculated pagination information
func (s *Service) ListFoldersByParentPaginated(userID string, parentID string, page, pageSize int) ([]Folder, int64, error) {
	// Ensure valid pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // Default page size
	}

	return s.repo.FindByUserAndParentPaginated(userID, parentID, page, pageSize)
}

// ListAllFoldersPaginated returns all folders for a user with pagination
// It also returns the total count of folders and calculated pagination information
func (s *Service) ListAllFoldersPaginated(userID string, page, pageSize int) ([]Folder, int64, error) {
	// Ensure valid pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // Default page size
	}

	return s.repo.FindAllByUserPaginated(userID, page, pageSize)
}

// GetFolderContents retrieves all files and folders within a specific folder
func (s *Service) GetFolderContents(folderID string, userID string) ([]Folder, []*file.File, error) {
	// Check if folder exists and belongs to user
	if folderID != "" {
		belongs, err := s.repo.BelongsToUser(folderID, userID)
		if err != nil {
			return nil, nil, err
		}
		if !belongs {
			return nil, nil, ErrFolderNotFound
		}
	}

	// Get subfolders
	folders, err := s.repo.FindByUserAndParent(userID, folderID)
	if err != nil {
		return nil, nil, err
	}

	// Get files in folder using the file repository
	files, err := s.fileService.ListFilesInFolder(userID, folderID)
	if err != nil {
		return nil, nil, err
	}

	return folders, files, nil
}

// DeleteFolder deletes a folder and all its contents (files and subfolders)
func (s *Service) DeleteFolder(folderID, userID string) error {
	// Check if folder exists and belongs to the user
	belongs, err := s.BelongsToUser(folderID, userID)
	if err != nil {
		return err
	}
	if !belongs {
		return ErrFolderNotFound
	}

	// Get all subfolders recursively
	subfolders, err := s.getAllSubfolders(userID, folderID)
	if err != nil {
		return err
	}

	// Delete all files in the folder and subfolders
	for _, subfolder := range append(subfolders, folderID) {
		if err := s.fileService.DeleteByFolder(userID, subfolder); err != nil {
			return err
		}
	}

	// Delete all subfolders
	for _, subfolder := range subfolders {
		if err := s.repo.Delete(subfolder); err != nil {
			return err
		}
	}

	// Delete the main folder
	return s.repo.Delete(folderID)
}

// getAllSubfolders recursively gets all subfolder IDs for a given folder
func (s *Service) getAllSubfolders(userID, folderID string) ([]string, error) {
	var result []string

	// Get direct children
	folders, err := s.repo.FindByUserAndParent(userID, folderID)
	if err != nil {
		return nil, err
	}

	// For each child folder
	for _, folder := range folders {
		// Add the folder ID to the result
		result = append(result, folder.ID)

		// Get all subfolders recursively
		subfolders, err := s.getAllSubfolders(userID, folder.ID)
		if err != nil {
			return nil, err
		}

		// Add all subfolders to the result
		result = append(result, subfolders...)
	}

	return result, nil
}
