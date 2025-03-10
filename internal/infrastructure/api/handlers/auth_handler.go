// internal/infrastructure/api/handlers/auth_handler.go
package handlers

import (
	"easy-storage/internal/domain/user"
	"easy-storage/internal/infrastructure/api/dto"
	"easy-storage/internal/infrastructure/auth/jwt"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication routes
type AuthHandler struct {
	userService *user.Service
	jwtProvider *jwt.Provider
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService *user.Service, jwtProvider *jwt.Provider) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtProvider: jwtProvider,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	// In a real app, add a validator here

	// Register the user
	newUser, err := h.userService.RegisterUser(req.Email, req.Password, req.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Generate tokens
	accessToken, err := h.jwtProvider.GenerateToken(newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	refreshToken, err := h.jwtProvider.GenerateRefreshToken(newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	// Return user info and tokens
	return c.Status(fiber.StatusCreated).JSON(dto.AuthResponse{
		User: dto.UserResponse{
			ID:    newUser.ID,
			Email: newUser.Email,
			Name:  newUser.Name,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    24 * 60 * 60, // 24 hours in seconds
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Authenticate user
	authenticatedUser, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		if err == user.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Authentication failed",
		})
	}

	// Generate tokens
	accessToken, err := h.jwtProvider.GenerateToken(authenticatedUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	refreshToken, err := h.jwtProvider.GenerateRefreshToken(authenticatedUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	// Return user info and tokens
	return c.Status(fiber.StatusOK).JSON(dto.AuthResponse{
		User: dto.UserResponse{
			ID:    authenticatedUser.ID,
			Email: authenticatedUser.Email,
			Name:  authenticatedUser.Name,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    24 * 60 * 60, // 24 hours in seconds
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate refresh token
	claims, err := h.jwtProvider.ValidateToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	// Get user from claims
	user, err := h.userService.GetUserByID(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Generate new tokens
	accessToken, err := h.jwtProvider.GenerateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	refreshToken, err := h.jwtProvider.GenerateRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	// Return new tokens
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    24 * 60 * 60, // 24 hours in seconds
	})
}

// GetMe retrieves the current user
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}
