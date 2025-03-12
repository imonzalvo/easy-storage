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

	// Get the showRootOnly parameter (optional, defaults to false)
	showRootOnly := c.QueryBool("showRootOnly", false)

	// Get pagination parameters from query
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Default page size with a reasonable limit
	}

	var folders []folder.Folder
	var totalCount int64
	var err error

	// Determine which folders to list based on parameters
	if parentID != "" {
		// If parentId is provided, show folders within that parent
		folders, totalCount, err = h.folderService.ListFoldersByParentPaginated(userID, parentID, page, pageSize)
	} else if showRootOnly {
		// If showRootOnly is true and no parentId, show only root folders
		folders, totalCount, err = h.folderService.ListFoldersByParentPaginated(userID, "", page, pageSize)
	} else {
		// Otherwise, show all folders
		folders, totalCount, err = h.folderService.ListAllFoldersPaginated(userID, page, pageSize)
	}

	if err != nil {
		log.Printf("error, %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not list folders",
		})
	}

	// Calculate total pages
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
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

	// Create pagination info
	paginationInfo := dto.PaginationInfo{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalItems:  totalCount,
		TotalPages:  totalPages,
		HasNextPage: page < totalPages,
		HasPrevPage: page > 1,
	}

	return c.Status(fiber.StatusOK).JSON(dto.FoldersListResponse{
		Folders:    folderResponses,
		Pagination: &paginationInfo,
	})
}

// GetFolderContents handles retrieving all contents of a folder
func (h *FolderHandler) GetFolderContents(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Get folder ID from parameter
	folderID := c.Params("folder_id")

	// Get folder contents
	folders, files, err := h.folderService.GetFolderContents(folderID, userID)
	if err != nil {
		if err == folder.ErrFolderNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Folder not found or you don't have permission to access it",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve folder contents",
		})
	}

	// Build response for folders
	folderResponses := make([]map[string]interface{}, len(folders))
	for i, folder := range folders {
		folderResponses[i] = map[string]interface{}{
			"id":         folder.ID,
			"name":       folder.Name,
			"parent_id":  folder.ParentID,
			"type":       "folder",
			"created_at": folder.CreatedAt.Format(time.RFC3339),
			"updated_at": folder.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Build response for files
	fileResponses := make([]map[string]interface{}, len(files))
	for i, file := range files {
		fileResponses[i] = map[string]interface{}{
			"id":           file.ID,
			"name":         file.Name,
			"size":         file.Size,
			"content_type": file.ContentType,
			"type":         "file",
			"created_at":   file.CreatedAt.Format(time.RFC3339),
			"updated_at":   file.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Combine both responses
	contents := append(folderResponses, fileResponses...)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"folder_id": folderID,
		"contents":  contents,
		"total":     len(contents),
	})
}

// DeleteFolder handles folder deletion
func (h *FolderHandler) DeleteFolder(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Get folder ID from parameter
	folderID := c.Params("folder_id")
	if folderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Folder ID is required",
		})
	}

	// Delete folder and all its contents
	err := h.folderService.DeleteFolder(folderID, userID)
	if err != nil {
		if err == folder.ErrFolderNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Folder not found or you don't have permission to delete it",
			})
		}
		log.Printf("Error deleting folder: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not delete folder",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Folder deleted successfully",
	})
}
