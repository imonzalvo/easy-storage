package common

// FolderValidator defines the interface for folder validation operations
type FolderValidator interface {
	// BelongsToUser checks if a folder belongs to a specific user
	BelongsToUser(folderID string, userID string) (bool, error)
}
