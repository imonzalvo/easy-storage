package folder

import "errors"

var (
	ErrFolderNotFound = errors.New("folder not found")
	ErrInvalidParent  = errors.New("invalid parent folder")
)

// Repository defines the interface for folder data access
type Repository interface {
	Save(folder *Folder) error
	FindByID(id string) (*Folder, error)
	FindByUserID(userID string) ([]*Folder, error)
	Delete(id string) error
	// Add method to check if folder belongs to user
	BelongsToUser(folderID string, userID string) (bool, error)
}
