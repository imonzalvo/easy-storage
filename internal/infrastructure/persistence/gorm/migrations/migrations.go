// internal/infrastructure/persistence/gorm/migrations/migrations.go
package migrations

import (
	"easy-storage/internal/infrastructure/persistence/gorm/models"

	"gorm.io/gorm"
)

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.Folder{},
	)
}
