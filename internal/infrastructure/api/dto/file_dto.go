// internal/infrastructure/api/dto/file_dto.go
package dto

// FileResponse represents file information returned to the client
type FileResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	FolderID    string `json:"folder_id,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// UploadFileResponse represents the response for a file upload
type UploadFileResponse struct {
	File FileResponse `json:"file"`
}

// FilesListResponse represents a list of files
type FilesListResponse struct {
	Files []FileResponse `json:"files"`
	Total int            `json:"total"`
}
