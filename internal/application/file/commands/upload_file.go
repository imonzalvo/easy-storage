// internal/application/file/commands/upload_file.go
package commands

import (
	"io"

	"easy-storage/internal/domain/file"
)

// UploadFileHandler handles file upload commands
type UploadFileHandler struct {
	fileService *file.Service
}

// NewUploadFileHandler creates a new upload file handler
func NewUploadFileHandler(fileService *file.Service) *UploadFileHandler {
	return &UploadFileHandler{
		fileService: fileService,
	}
}

// UploadFileCommand represents a command to upload a file
type UploadFileCommand struct {
	Filename    string
	Size        int64
	ContentType string
	Content     io.Reader
	UserID      string
	FolderID    string
}

// Handle executes the upload file command
func (h *UploadFileHandler) Handle(cmd *UploadFileCommand) (*file.File, error) {
	return h.fileService.UploadFile(
		cmd.Filename,
		cmd.Size,
		cmd.ContentType,
		cmd.Content,
		cmd.UserID,
		cmd.FolderID,
	)
}
