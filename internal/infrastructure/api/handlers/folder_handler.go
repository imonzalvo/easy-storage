package handlers

import (
	"easy-storage/internal/domain/folder"
	"easy-storage/internal/infrastructure/api/dto"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// FolderHandler handles folder-related API endpoints
type FolderHandler struct {
	folderService *folder.Service
}

// NewFolderHandler creates a new folder handler
func NewFolderHandler(folderService *folder.Service) *FolderHandler {
	return &FolderHandler{
		folderService: folderService,
	}
}

// CreateFolder handles folder creation
func (h *FolderHandler) CreateFolder(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Parse request body
	var req dto.CreateFolderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate folder name
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Folder name is required",
		})
	}

	// Create folder
	createdFolder, err := h.folderService.CreateFolder(req.Name, req.ParentID, userID)
	if err != nil {
		if err == folder.ErrInvalidParent {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid parent folder",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create folder",
		})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(dto.FolderResponse{
		ID:        createdFolder.ID,
		Name:      createdFolder.Name,
		ParentID:  createdFolder.ParentID,
		CreatedAt: createdFolder.CreatedAt.Format(time.RFC3339),
		UpdatedAt: createdFolder.UpdatedAt.Format(time.RFC3339),
	})
}

// ListFolders handles listing folders for the current user within a specific parent folder
func (h *FolderHandler) ListFolders(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Get parent folder ID from query parameter (optional)
	parentID := c.Query("parentId", "")

	// List folders
	folders, err := h.folderService.ListFoldersByParent(userID, parentID)
	if err != nil {
		log.Printf("error, %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not list folders",
		})
	}

	// Build response
	folderResponses := make([]dto.FolderResponse, len(folders))
	for i, folder := range folders {
		folderResponses[i] = dto.FolderResponse{
			ID:        folder.ID,
			Name:      folder.Name,
			ParentID:  folder.ParentID,
			CreatedAt: folder.CreatedAt.Format(time.RFC3339),
			UpdatedAt: folder.UpdatedAt.Format(time.RFC3339),
		}
	}

	return c.Status(fiber.StatusOK).JSON(dto.FoldersListResponse{
		Folders: folderResponses,
	})
}
