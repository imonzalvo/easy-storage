// internal/infrastructure/api/router.go
package api

import (
	"easy-storage/internal/domain/file"
	"easy-storage/internal/domain/folder"
	"easy-storage/internal/domain/user"
	"easy-storage/internal/infrastructure/api/handlers"
	"easy-storage/internal/infrastructure/api/middleware"
	"easy-storage/internal/infrastructure/auth/jwt"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	app *fiber.App,
	userService *user.Service,
	fileService *file.Service,
	folderService *folder.Service,
	jwtProvider *jwt.Provider,
) {
	authHandler := handlers.NewAuthHandler(userService, jwtProvider)
	fileHandler := handlers.NewFileHandler(fileService)
	folderHandler := handlers.NewFolderHandler(folderService)

	// Auth routes
	auth := app.Group("/api/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	api := app.Group("/api", middleware.AuthMiddleware(jwtProvider))
	api.Get("/me", authHandler.GetMe)

	// File routes
	fileRoutes := api.Group("/files")
	fileRoutes.Post("/", fileHandler.UploadFile)
	fileRoutes.Get("/", fileHandler.ListFiles)
	fileRoutes.Get("/:id", fileHandler.DownloadFile)
	fileRoutes.Delete("/:id", fileHandler.DeleteFile)

	// Folder routes
	folderRoutes := api.Group("/folders")
	folderRoutes.Post("/", folderHandler.CreateFolder)
	// folderRoutes.Get("/", folderHandler.ListFolders)
	// folderRoutes.Get("/:id", folderHandler.GetFolder)
	// folderRoutes.Delete("/:id", folderHandler.DeleteFolder)
}
