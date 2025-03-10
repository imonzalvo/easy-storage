package folder

// Repository is used here to implement the common.FolderValidator interface
// This allows the folder package to provide validation functionality without
// creating a cyclic dependency

// BelongsToUser checks if a folder belongs to a specific user
func (s *Service) BelongsToUser(folderID string, userID string) (bool, error) {
	return s.repo.BelongsToUser(folderID, userID)
}
