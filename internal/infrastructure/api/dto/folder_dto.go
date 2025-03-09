package dto

// CreateFolderRequest represents the request to create a folder
type CreateFolderRequest struct {
	Name     string `json:"name"`
	ParentID string `json:"parent_id,omitempty"`
}

// FolderResponse represents folder information returned to the client
type FolderResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ParentID  string `json:"parent_id,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// FoldersListResponse represents a list of folders
type FoldersListResponse struct {
	Folders []FolderResponse `json:"folders"`
}
