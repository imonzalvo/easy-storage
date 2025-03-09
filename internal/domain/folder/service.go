package folder

// Service provides folder operations
type Service struct {
	repo Repository
}

// NewService creates a new folder service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
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

// DeleteFolder deletes a folder
func (s *Service) DeleteFolder(id string) error {
	return s.repo.Delete(id)
}
