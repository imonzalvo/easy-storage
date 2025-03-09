// internal/infrastructure/api/middleware/auth.go
package middleware

import (
	"strings"

	"easy-storage/internal/infrastructure/auth/jwt"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates middleware for JWT authentication
func AuthMiddleware(jwtProvider *jwt.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is missing",
			})
		}

		// Check if the header has the Bearer prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		// Extract the token
		tokenString := tokenParts[1]

		// Validate the token
		claims, err := jwtProvider.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Set user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}
