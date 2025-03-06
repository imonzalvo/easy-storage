package commands

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"easy-storage/internal/domain/file"
)

// UploadFileCommand represents the command to upload a file
type UploadFileCommand struct {
	Name        string
	FileHeader  *multipart.FileHeader
	ContentType string
	UserID      uuid.UUID
	FolderID    *uuid.UUID
}

// UploadFileHandler handles the upload file command
type UploadFileHandler struct {
	fileService *file.Service
}

// NewUploadFileHandler creates a new upload file handler
func NewUploadFileHandler(fileService *file.Service) *UploadFileHandler {
	return &UploadFileHandler{
		fileService: fileService,
	}
}

// Handle processes the upload file command
func (h *UploadFileHandler) Handle(ctx context.Context, cmd UploadFileCommand) (*file.File, error) {
	return h.fileService.UploadFile(
		ctx,
		cmd.Name,
		cmd.FileHeader,
		cmd.ContentType,
		cmd.UserID,
		cmd.FolderID,
	)
}
