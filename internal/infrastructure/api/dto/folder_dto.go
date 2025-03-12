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

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	HasNextPage bool  `json:"has_next_page"`
	HasPrevPage bool  `json:"has_prev_page"`
}

// FoldersListResponse represents a list of folders
type FoldersListResponse struct {
	Folders    []FolderResponse `json:"folders"`
	Pagination *PaginationInfo  `json:"pagination,omitempty"`
}
