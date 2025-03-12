package handlers

import (
	"fmt"
	"time"

	"easy-storage/internal/domain/access"
	"easy-storage/internal/domain/file"
	"easy-storage/internal/domain/user"

	"github.com/gofiber/fiber/v2"
)

// FileHandler handles file-related API endpoints
type FileHandler struct {
	fileService   *file.Service
	accessService *access.Service
}

// NewFileHandler creates a new file handler
func NewFileHandler(fileService *file.Service, accessService *access.Service) *FileHandler {
	return &FileHandler{
		fileService:   fileService,
		accessService: accessService,
	}
}

// UploadFile handles file uploads
func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Get folder ID from query parameter (optional)
	folderID := c.Query("folder_id", "")

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file provided",
		})
	}

	// Check file size (example: limit to 100MB)
	if file.Size > 100*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File too large, maximum size is 100MB",
		})
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not open uploaded file",
		})
	}
	defer src.Close()

	// Determine content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload file
	uploadedFile, err := h.fileService.UploadFile(
		file.Filename,
		file.Size,
		contentType,
		src,
		userID,
		folderID,
	)
	if err != nil {
		if err == user.ErrStorageQuotaExceeded {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Storage quota exceeded",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Could not upload file: %v", err),
		})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":           uploadedFile.ID,
		"name":         uploadedFile.Name,
		"size":         uploadedFile.Size,
		"content_type": uploadedFile.ContentType,
		"folder_id":    uploadedFile.FolderID,
		"created_at":   uploadedFile.CreatedAt.Format(time.RFC3339),
		"updated_at":   uploadedFile.UpdatedAt.Format(time.RFC3339),
	})
}

// DownloadFile now returns a signed URL for file download
func (h *FileHandler) DownloadFile(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Get file ID from parameter
	fileID := c.Params("id")
	if fileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File ID is required",
		})
	}

	// Check if file belongs to user or has been shared with user
	hasAccess, err := h.accessService.CheckFileAccess(c.Context(), fileID, userID)
	if err != nil || !hasAccess {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to access this file",
		})
	}

	// Get file from service
	downloadedFile, err := h.fileService.GetFile(fileID)
	if err != nil {
		if err == file.ErrFileNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "File not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve file",
		})
	}

	// Get signed URL (valid for 1 hour = 3600 seconds)
	signedURL, err := h.fileService.GetFileSignedURL(downloadedFile, 3600)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate download URL",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"url":          signedURL,
		"expires_in":   3600,
		"filename":     downloadedFile.Name,
		"content_type": downloadedFile.ContentType,
		"size":         downloadedFile.Size,
	})
}

// ListFiles handles listing files for the current user
func (h *FileHandler) ListFiles(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Get pagination parameters from query
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	// Get sort parameter (default to created_at)
	sort := c.Query("sort", "created_at")

	// Get sort direction (default to desc)
	sortDir := c.Query("sort_dir", "desc")

	// Validate sort parameter
	validSorts := map[string]bool{
		"name":       true,
		"size":       true,
		"created_at": true,
	}

	if !validSorts[sort] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sort parameter. Valid values are: name, size, created_at",
		})
	}

	// Validate sort direction
	if sortDir != "asc" && sortDir != "desc" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sort_dir parameter. Valid values are: asc, desc",
		})
	}

	// List files
	files, err := h.fileService.ListUserFiles(userID, limit, offset, sort, sortDir)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not list files",
		})
	}

	// Build response
	fileResponses := make([]map[string]interface{}, len(files))
	for i, file := range files {
		fileResponses[i] = map[string]interface{}{
			"id":           file.ID,
			"name":         file.Name,
			"size":         file.Size,
			"content_type": file.ContentType,
			"folder_id":    file.FolderID,
			"created_at":   file.CreatedAt.Format(time.RFC3339),
			"updated_at":   file.UpdatedAt.Format(time.RFC3339),
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"files": fileResponses,
		"total": len(files),
	})
}

// DeleteFile handles file deletion
func (h *FileHandler) DeleteFile(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Get file ID from parameter
	fileID := c.Params("id")
	if fileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File ID is required",
		})
	}

	// Get file from service to check ownership
	deletedFile, err := h.fileService.GetFile(fileID)
	if err != nil {
		if err == file.ErrFileNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "File not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve file",
		})
	}

	// Check if file belongs to user
	if deletedFile.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to delete this file",
		})
	}

	// Delete file
	if err := h.fileService.DeleteFile(fileID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not delete file",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
