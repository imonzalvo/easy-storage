package main

import (
	"log"

	"easy-storage/internal/config"
	"easy-storage/internal/domain/access"
	"easy-storage/internal/domain/file"
	"easy-storage/internal/domain/folder"
	"easy-storage/internal/domain/share"
	"easy-storage/internal/domain/user"
	"easy-storage/internal/infrastructure/api"
	"easy-storage/internal/infrastructure/auth/jwt"
	"easy-storage/internal/infrastructure/persistence"
	"easy-storage/internal/infrastructure/persistence/gorm/repositories"
	"easy-storage/internal/infrastructure/storage/s3"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := persistence.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize storage provider
	storageProvider, err := s3.NewS3Provider(&cfg.Storage)
	if err != nil {
		log.Fatalf("Failed to initialize storage provider: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewGormUserRepository(db)
	fileRepo := repositories.NewGormFileRepository(db)
	folderRepo := repositories.NewGormFolderRepository(db)
	shareRepo := repositories.NewShareRepository(db) // Add share repository

	// Initialize domain services
	userService := user.NewService(userRepo)
	fileService := file.NewService(fileRepo, folderRepo, storageProvider)
	folderService := folder.NewService(folderRepo, fileService)
	shareService := share.NewService(shareRepo)
	accessService := access.NewService(fileService, shareService)

	// Initialize JWT provider
	jwtProvider := jwt.NewProvider(
		cfg.Auth.JWTSecret,
		cfg.Auth.TokenExpiry,
		cfg.Auth.RefreshExpiry,
	)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "easy-storage",
	})

	// Middleware
	app.Use(fiberLogger.New())
	app.Use(cors.New())

	// Setup routes
	api.SetupRoutes(app, userService, fileService, folderService, shareService, accessService, jwtProvider)

	// Default route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("easy-storage API is running")
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
