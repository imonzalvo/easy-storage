package handlers

import (
	"log"
	"time"

	"easy-storage/internal/domain/access"
	"easy-storage/internal/domain/file"
	"easy-storage/internal/domain/share"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ShareHandler handles share-related API endpoints
type ShareHandler struct {
	shareService      *share.Service
	fileService       *file.Service
	fileAccessService *access.Service
}

// NewShareHandler creates a new share handler
func NewShareHandler(shareService *share.Service, fileService *file.Service, fileAccessService *access.Service) *ShareHandler {
	return &ShareHandler{
		shareService:      shareService,
		fileService:       fileService,
		fileAccessService: fileAccessService,
	}
}

// CreateShare handles creating a new share
func (h *ShareHandler) CreateShare(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Parse request body
	var req struct {
		ResourceID   string `json:"resource_id" validate:"required"`
		ResourceType string `json:"resource_type" validate:"required,oneof=file folder"`
		ShareType    string `json:"share_type" validate:"required,oneof=LINK USER"`
		Permission   string `json:"permission" validate:"required,oneof=READ WRITE"`
		RecipientID  string `json:"recipient_id,omitempty"`
		Password     string `json:"password,omitempty"`
		ExpiresAt    string `json:"expires_at,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Convert string IDs to UUID
	ownerID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	resourceID, err := uuid.Parse(req.ResourceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid resource ID",
		})
	}

	// Create share based on type
	var newShare *share.Share
	if req.ShareType == "LINK" {
		// Create link share
		newShare, err = h.shareService.CreateLinkShare(
			c.Context(),
			ownerID,
			resourceID,
			req.ResourceType,
			share.SharePermission(req.Permission),
		)
	} else {
		// For user shares, recipient is required
		if req.RecipientID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Recipient ID is required for user shares",
			})
		}

		recipientID, err := uuid.Parse(req.RecipientID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid recipient ID",
			})
		}

		// Create user share
		newShare, err = h.shareService.CreateUserShare(
			c.Context(),
			ownerID,
			resourceID,
			req.ResourceType,
			recipientID,
			share.SharePermission(req.Permission),
		)
	}

	if err != nil {
		log.Printf("error %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create share",
		})
	}

	// Set optional parameters if provided
	if req.Password != "" {
		if err := h.shareService.SetSharePassword(c.Context(), newShare.ID, req.Password); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not set password",
			})
		}
	}

	if req.ExpiresAt != "" {
		expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid expiration date format, use RFC3339",
			})
		}

		if err := h.shareService.SetShareExpiration(c.Context(), newShare.ID, expiresAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not set expiration",
			})
		}
	}

	// Get the updated share
	newShare, err = h.shareService.GetShareByID(c.Context(), newShare.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve created share",
		})
	}

	// Build response
	response := buildShareResponse(newShare)
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetShare handles retrieving a share by ID
func (h *ShareHandler) GetShare(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Get share ID from parameter
	shareID := c.Params("id")
	if shareID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Share ID is required",
		})
	}

	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	shareUUID, err := uuid.Parse(shareID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid share ID",
		})
	}

	// Get share
	existingShare, err := h.shareService.GetShareByID(c.Context(), shareUUID)
	if err != nil {
		if err == share.ErrShareNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Share not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve share",
		})
	}

	// Check if user is authorized to view this share
	isOwner := existingShare.OwnerID == userUUID
	isRecipient := existingShare.RecipientID != nil && *existingShare.RecipientID == userUUID

	if !isOwner && !isRecipient {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to access this share",
		})
	}

	// Build response
	response := buildShareResponse(existingShare)
	return c.Status(fiber.StatusOK).JSON(response)
}

// ListShares handles listing shares for the current user
func (h *ShareHandler) ListShares(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// List shares
	shares, err := h.shareService.ListSharesByOwner(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not list shares",
		})
	}

	// Build response
	shareResponses := make([]map[string]interface{}, len(shares))
	for i, s := range shares {
		shareResponses[i] = buildShareResponse(s)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"shares": shareResponses,
		"total":  len(shares),
	})
}

// ListSharesWithMe handles listing shares shared with the current user
func (h *ShareHandler) ListSharesWithMe(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// List shares
	shares, err := h.shareService.ListSharesWithUser(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not list shares",
		})
	}

	// Build response
	shareResponses := make([]map[string]interface{}, len(shares))
	for i, s := range shares {
		shareResponses[i] = buildShareResponse(s)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"shares": shareResponses,
		"total":  len(shares),
	})
}

// ListSharesByResource handles listing all shares for a specific resource
func (h *ShareHandler) ListSharesByResource(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Get resource type and ID from parameters
	resourceType := c.Params("type")
	resourceID := c.Params("id")

	if resourceID == "" || resourceType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource ID and type are required",
		})
	}

	// Validate resource type
	if resourceType != "file" && resourceType != "folder" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid resource type, must be 'file' or 'folder'",
		})
	}

	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	resourceUUID, err := uuid.Parse(resourceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid resource ID",
		})
	}

	// List shares for the resource
	shares, err := h.shareService.ListSharesByResource(c.Context(), resourceUUID, resourceType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not list shares",
		})
	}

	// Filter to only include shares owned by the requesting user
	ownedShares := make([]*share.Share, 0)
	for _, s := range shares {
		if s.OwnerID == userUUID {
			ownedShares = append(ownedShares, s)
		}
	}

	// Build response
	shareResponses := make([]map[string]interface{}, len(ownedShares))
	for i, s := range ownedShares {
		shareResponses[i] = buildShareResponse(s)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"shares": shareResponses,
		"total":  len(ownedShares),
	})
}

// RevokeShare handles revoking a share
func (h *ShareHandler) RevokeShare(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Get share ID from parameter
	shareID := c.Params("id")
	if shareID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Share ID is required",
		})
	}

	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	shareUUID, err := uuid.Parse(shareID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid share ID",
		})
	}

	// Get share to check ownership
	existingShare, err := h.shareService.GetShareByID(c.Context(), shareUUID)
	if err != nil {
		if err == share.ErrShareNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Share not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve share",
		})
	}

	// Check if user is the owner
	if existingShare.OwnerID != userUUID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to revoke this share",
		})
	}

	// Revoke share
	if err := h.shareService.RevokeShare(c.Context(), shareUUID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not revoke share",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateShare handles updating a share's properties
func (h *ShareHandler) UpdateShare(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Get share ID from parameter
	shareID := c.Params("id")
	if shareID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Share ID is required",
		})
	}

	// Parse request body
	var req struct {
		Password  string `json:"password"`
		ExpiresAt string `json:"expires_at"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	shareUUID, err := uuid.Parse(shareID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid share ID",
		})
	}

	// Get share to check ownership
	existingShare, err := h.shareService.GetShareByID(c.Context(), shareUUID)
	if err != nil {
		if err == share.ErrShareNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Share not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve share",
		})
	}

	// Check if user is the owner
	if existingShare.OwnerID != userUUID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this share",
		})
	}

	// Update password if provided
	if req.Password != "" {
		if err := h.shareService.SetSharePassword(c.Context(), shareUUID, req.Password); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not update password",
			})
		}
	}

	// Update expiration if provided
	if req.ExpiresAt != "" {
		// Parse expiration time
		expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid expiration date format, use RFC3339",
			})
		}

		if err := h.shareService.SetShareExpiration(c.Context(), shareUUID, expiresAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not update expiration",
			})
		}
	}

	// Get updated share
	updatedShare, err := h.shareService.GetShareByID(c.Context(), shareUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve updated share",
		})
	}

	// Build response
	response := buildShareResponse(updatedShare)
	return c.Status(fiber.StatusOK).JSON(response)
}

// AccessShare handles accessing a shared resource via token
func (h *ShareHandler) AccessShare(c *fiber.Ctx) error {
	// Get token from parameter
	token := c.Params("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is required",
		})
	}

	// Parse password from query if provided
	password := c.Query("password")

	// Get share by token
	existingShare, err := h.shareService.GetShareByToken(c.Context(), token)
	if err != nil {
		if err == share.ErrShareNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Share not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve share",
		})
	}

	// Check if share is accessible
	if !existingShare.IsAccessible() {
		if existingShare.IsRevoked {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "This share has been revoked",
			})
		}
		if existingShare.IsExpired() {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "This share has expired",
			})
		}
	}

	// Check password if required
	if existingShare.Password != nil {
		// If no password provided but required
		if password == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":             "Password required",
				"requires_password": true,
			})
		}

		// If password doesn't match
		if *existingShare.Password != password {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid password",
			})
		}
	}

	// Record access
	if err := h.shareService.RecordShareAccess(c.Context(), existingShare.ID); err != nil {
		// Log this error but don't fail the request
		// logger.Warn("Failed to record share access", "error", err)
	}

	// Build response
	response := buildShareResponse(existingShare)
	return c.Status(fiber.StatusOK).JSON(response)
}

// ValidateShareAccess handles validating access to a shared resource with password
func (h *ShareHandler) ValidateShareAccess(c *fiber.Ctx) error {
	// Get token from parameter
	token := c.Params("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is required",
		})
	}

	// Parse request body
	var req struct {
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get share by token
	existingShare, err := h.shareService.GetShareByToken(c.Context(), token)
	if err != nil {
		if err == share.ErrShareNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Share not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve share",
		})
	}

	// Check if share is accessible
	if !existingShare.IsAccessible() {
		if existingShare.IsRevoked {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "This share has been revoked",
			})
		}
		if existingShare.IsExpired() {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "This share has expired",
			})
		}
	}

	// Check password
	if existingShare.Password != nil {
		if req.Password == "" || *existingShare.Password != req.Password {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid password",
			})
		}
	}

	// Record access
	if err := h.shareService.RecordShareAccess(c.Context(), existingShare.ID); err != nil {
		// Log this error but don't fail the request
		// logger.Warn("Failed to record share access", "error", err)
	}

	// Build response
	response := buildShareResponse(existingShare)
	return c.Status(fiber.StatusOK).JSON(response)
}

// Helper function to build share response
func buildShareResponse(s *share.Share) map[string]interface{} {
	response := map[string]interface{}{
		"id":            s.ID,
		"resource_id":   s.ResourceID,
		"resource_type": s.ResourceType,
		"share_type":    s.Type,
		"permission":    s.Permission,
		"has_password":  s.Password != nil,
		"is_revoked":    s.IsRevoked,
		"access_count":  s.AccessCount,
		"created_at":    s.CreatedAt.Format(time.RFC3339),
		"updated_at":    s.UpdatedAt.Format(time.RFC3339),
	}

	// Add optional fields if present
	if s.ExpiresAt != nil {
		response["expires_at"] = s.ExpiresAt.Format(time.RFC3339)
	}

	if s.LastAccessAt != nil {
		response["last_access_at"] = s.LastAccessAt.Format(time.RFC3339)
	}

	if s.RecipientID != nil {
		response["recipient_id"] = s.RecipientID
	}

	// Only include token in response for link shares
	if s.Type == share.LinkShare {
		response["token"] = s.Token
		// You could add a full URL here if you have a base URL configured
		// response["url"] = fmt.Sprintf("%s/share/%s", baseURL, s.Token)
	}

	return response
}

// DownloadSharedFile handles downloading a file via share token
func (h *ShareHandler) DownloadSharedFile(c *fiber.Ctx) error {
	// Get token from parameter
	token := c.Params("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is required",
		})
	}

	// Parse password from query if provided
	password := c.Query("password", "")

	// Get file using share token
	downloadedFile, err := h.fileAccessService.GetFileByShareToken(c.Context(), token, password)
	if err != nil {
		switch err {
		case share.ErrShareNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Share not found",
			})
		case share.ErrShareRevoked:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "This share has been revoked",
			})
		case share.ErrShareExpired:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "This share has expired",
			})
		case share.ErrInvalidPassword:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":             "Invalid password",
				"requires_password": true,
			})
		case access.ErrInvalidResourceType:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Only files can be downloaded",
			})
		case file.ErrFileNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "File not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not retrieve file",
			})
		}
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
