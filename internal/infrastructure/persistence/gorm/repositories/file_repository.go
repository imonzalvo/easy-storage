package repositories

import (
	"errors"

	"easy-storage/internal/domain/file"
	"easy-storage/internal/infrastructure/persistence/gorm/models"

	"gorm.io/gorm"
)

// GormFileRepository implements the file.Repository interface using GORM
type GormFileRepository struct {
	db *gorm.DB
}

// NewGormFileRepository creates a new file repository
func NewGormFileRepository(db *gorm.DB) file.Repository {
	return &GormFileRepository{db: db}
}

// Save creates or updates a file in the database
func (r *GormFileRepository) Save(f *file.File) error {
	fileModel := &models.File{
		ID:          f.ID,
		Name:        f.Name,
		Size:        f.Size,
		ContentType: f.ContentType,
		Path:        f.Path,
		UserID:      f.UserID,
		FolderID:    f.FolderID,
	}

	if err := r.db.Save(fileModel).Error; err != nil {
		return err
	}

	f.ID = fileModel.ID
	return nil
}

// FindByID finds a file by ID
func (r *GormFileRepository) FindByID(id string) (*file.File, error) {
	var fileModel models.File
	if err := r.db.First(&fileModel, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, file.ErrFileNotFound
		}
		return nil, err
	}

	return &file.File{
		ID:          fileModel.ID,
		Name:        fileModel.Name,
		Size:        fileModel.Size,
		ContentType: fileModel.ContentType,
		Path:        fileModel.Path,
		UserID:      fileModel.UserID,
		FolderID:    fileModel.FolderID,
		CreatedAt:   fileModel.CreatedAt,
		UpdatedAt:   fileModel.UpdatedAt,
	}, nil
}

// FindByUserID finds files by user ID with pagination
func (r *GormFileRepository) FindByUserID(userID string, limit, offset int) ([]*file.File, error) {
	var fileModels []models.File
	if err := r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&fileModels).Error; err != nil {
		return nil, err
	}

	files := make([]*file.File, len(fileModels))
	for i, fileModel := range fileModels {
		files[i] = &file.File{
			ID:          fileModel.ID,
			Name:        fileModel.Name,
			Size:        fileModel.Size,
			ContentType: fileModel.ContentType,
			Path:        fileModel.Path,
			UserID:      fileModel.UserID,
			FolderID:    fileModel.FolderID,
			CreatedAt:   fileModel.CreatedAt,
			UpdatedAt:   fileModel.UpdatedAt,
		}
	}

	return files, nil
}

// Delete deletes a file
func (r *GormFileRepository) Delete(id string) error {
	return r.db.Delete(&models.File{}, "id = ?", id).Error
}

// FindByUserIDAndFolder finds files by user ID and folder ID with pagination
func (r *GormFileRepository) FindByUserIDAndFolder(userID string, folderID string) ([]*file.File, error) {
	var fileModels []models.File
	query := r.db.Where("user_id = ?", userID)

	if folderID == "" {
		// Find files in the root folder (where folder_id is null or empty)
		query = query.Where("folder_id IS NULL OR folder_id = ''")
	} else {
		// Find files in the specified folder
		query = query.Where("folder_id = ?", folderID)
	}

	if err := query.Find(&fileModels).Error; err != nil {
		return nil, err
	}

	files := make([]*file.File, len(fileModels))
	for i, fileModel := range fileModels {
		files[i] = &file.File{
			ID:          fileModel.ID,
			Name:        fileModel.Name,
			Size:        fileModel.Size,
			ContentType: fileModel.ContentType,
			Path:        fileModel.Path,
			UserID:      fileModel.UserID,
			FolderID:    fileModel.FolderID,
			CreatedAt:   fileModel.CreatedAt,
			UpdatedAt:   fileModel.UpdatedAt,
		}
	}

	return files, nil
}

// DeleteByFolder deletes all files in the database that belong to a specific folder
func (r *GormFileRepository) DeleteByFolder(folderID string) error {
	// Find all files in the folder to get their paths
	var files []models.File
	if err := r.db.Where("folder_id = ?", folderID).Find(&files).Error; err != nil {
		return err
	}

	// We don't need to delete from storage here since that's handled by the file service

	// Delete all files from database
	return r.db.Delete(&models.File{}, "folder_id = ?", folderID).Error
}
