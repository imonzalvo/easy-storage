// internal/infrastructure/persistence/db.go
package persistence

import (
	"fmt"
	"log"
	"time"

	"easy-storage/internal/config"
	"easy-storage/internal/infrastructure/persistence/gorm/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase creates a new database connection
func NewDatabase(config *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name,
	)

	// Retry connection a few times
	var db *gorm.DB
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to connect to database (attempt %d/%d)...", i+1, maxRetries)
		log.Printf("DSN (attempt %s)...", dsn)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			break
		}

		log.Printf("Failed to connect to database: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	log.Println("Successfully connected to database")

	// Run migrations
	if err := migrations.Migrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}
