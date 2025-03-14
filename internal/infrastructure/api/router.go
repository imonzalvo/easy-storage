// internal/infrastructure/api/router.go
package api

import (
	"easy-storage/internal/domain/access"
	"easy-storage/internal/domain/file"
	"easy-storage/internal/domain/folder"
	"easy-storage/internal/domain/share"
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
	shareService *share.Service,
	accessService *access.Service,
	jwtProvider *jwt.Provider,
) {
	authHandler := handlers.NewAuthHandler(userService, jwtProvider)
	fileHandler := handlers.NewFileHandler(fileService, accessService)
	folderHandler := handlers.NewFolderHandler(folderService)
	shareHandler := handlers.NewShareHandler(shareService, fileService, accessService)

	// Auth routes
	auth := app.Group("/api/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	api := app.Group("/api", middleware.AuthMiddleware(jwtProvider))
	api.Get("/me", authHandler.GetMe)
	api.Post("/auth/change-password", authHandler.ChangePassword)

	// File routes
	fileRoutes := api.Group("/files")
	fileRoutes.Post("/", fileHandler.UploadFile)
	fileRoutes.Get("/", fileHandler.ListFiles)
	fileRoutes.Get("/:id", fileHandler.DownloadFile)
	fileRoutes.Delete("/:id", fileHandler.DeleteFile)

	// Folder routes
	folderRoutes := api.Group("/folders")
	folderRoutes.Post("/", folderHandler.CreateFolder)
	folderRoutes.Get("/", folderHandler.ListFolders)
	folderRoutes.Get("/:folder_id", folderHandler.GetFolderContents)
	folderRoutes.Delete("/:folder_id", folderHandler.DeleteFolder)

	// Share routes
	shareGroup := app.Group("/api/shares")
	shareGroup.Post("/", shareHandler.CreateShare)
	shareGroup.Get("/", shareHandler.ListShares)
	shareGroup.Get("/shared-with-me", shareHandler.ListSharesWithMe)
	shareGroup.Get("/:id", shareHandler.GetShare)
	shareGroup.Delete("/:id", shareHandler.RevokeShare)

	// Public share access endpoint (no auth required)
	app.Get("/share/:token", shareHandler.AccessShare)
	app.Get("/share/:token/download", shareHandler.DownloadSharedFile)
}
