package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"easy-storage/internal/application/file/commands"
	"easy-storage/internal/application/file/queries"
	"easy-storage/internal/domain/file"
	"easy-storage/internal/infrastructure/api/dto"
)

// FileHandler handles HTTP requests related to files
type FileHandler struct {
	uploadFileHandler *commands.UploadFileHandler
	// deleteFileHandler *commands.DeleteFileHandler
	// getFileHandler    *queries.GetFileHandler
	// listFilesHandler  *queries.ListFilesHandler
}

// NewFileHandler creates a new file handler
func NewFileHandler(
	uploadFileHandler *commands.UploadFileHandler,
	// deleteFileHandler *commands.DeleteFileHandler,
	// getFileHandler *queries.GetFileHandler,
	// listFilesHandler *queries.ListFilesHandler,
) *FileHandler {
	return &FileHandler{
		uploadFileHandler: uploadFileHandler,
		// deleteFileHandler: deleteFileHandler,
		// getFileHandler:    getFileHandler,
		// listFilesHandler:  listFilesHandler,
	}
}

// RegisterRoutes registers the file routes
func (h *FileHandler) RegisterRoutes(app *fiber.App) {
	files := app.Group("/api/files")
	
	files.Post("/", h.UploadFile)
	// files.Get("/", h.ListFiles)
	// files.Get("/:id", h.GetFile)
	// files.Delete("/:id", h.DeleteFile)
}

// UploadFile handles file upload
func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	// Get the authenticated user ID from context
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	
	// Parse multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file upload",
		})
	}
	
	// Get folder ID if provided
	var folderID *uuid.UUID
	if folderIDStr := c.FormValue("folder_id"); folderIDStr != "" {
		id, err := uuid.Parse(folderIDStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid folder ID",
			})
		}
		folderID = &id
	}
	
	// Get content type
	contentType := file.Header.Get("Content-Type")
	
	// Create command
	cmd := commands.UploadFileCommand{
		Name:        file.Filename,
		FileHeader:  file,
		ContentType: contentType,
		UserID:      userID,
		FolderID:    folderID,
	}
	
	// Handle command
	result, err := h.uploadFileHandler.Handle(c.Context(), cmd)
	if err != nil {
		// Handle specific domain errors
		switch err {
		case file.ErrStorageQuotaExceeded:
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "Storage quota exceeded",
			})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to upload file: %v", err),
			})
		}
	}
	
	// Convert to DTO
	fileDTO := dto.FileToDTO(result)
	
	return c.Status(http.StatusCreated).JSON(fileDTO)
}

