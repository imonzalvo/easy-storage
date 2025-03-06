package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/imonzalvo/easy-storage/internal/config"
	"github.com/imonzalvo/easy-storage/internal/infrastructure/api/handlers"
	"github.com/imonzalvo/easy-storage/internal/infrastructure/persistence/gorm"
	"github.com/imonzalvo/easy-storage/internal/infrastructure/storage/s3"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := gorm.NewConnection(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := gorm.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize storage
	storage, err := s3.NewS3Storage(
		cfg.Storage.Bucket,
		cfg.Storage.Region,
		cfg.Storage.Endpoint,
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
	)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize repositories
	userRepo := gorm.NewUserRepository(db)
	fileRepo := gorm.NewFileRepository(db)
	folderRepo := gorm.NewFolderRepository(db)
	shareRepo := gorm.NewShareRepository(db)

	// Initialize domain services
	fileService := file.NewService(fileRepo, storage, userRepo)
	folderService := folder.NewService(folderRepo, fileRepo)
	userService := user.NewService(userRepo)
	shareService := share.NewService(shareRepo, fileRepo, folderRepo)

	// Initialize application handlers
	uploadFileHandler := commands.NewUploadFileHandler(fileService)
	getFileHandler := queries.NewGetFileHandler(fileService)
	// ... other handlers initialization

	// Initialize API handlers
	fileHandler := handlers.NewFileHandler(uploadFileHandler, deleteFileHandler, getFileHandler, listFilesHandler)
	// ... other API handlers initialization

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "easy-storage",
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Register routes
	fileHandler.RegisterRoutes(app)
	// ... other route registrations

	// Start server
	go func() {
		if err := app.Listen(cfg.Server.Address); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	// Custom error handling logic
	code := fiber.StatusInternalServerError

	// Check for specific error types
	// ... error type checking logic

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
