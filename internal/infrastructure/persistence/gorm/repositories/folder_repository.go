package repositories

import (
	"easy-storage/internal/domain/folder"
	"easy-storage/internal/infrastructure/persistence/gorm/models"
	"errors"

	"gorm.io/gorm"
)

// GormFolderRepository implements the folder.Repository interface using GORM
type GormFolderRepository struct {
	db *gorm.DB
}

// NewGormFolderRepository creates a new folder repository
func NewGormFolderRepository(db *gorm.DB) folder.Repository {
	return &GormFolderRepository{db: db}
}

// Save creates or updates a folder in the database
func (r *GormFolderRepository) Save(f *folder.Folder) error {
	folderModel := &models.Folder{
		ID:       f.ID,
		Name:     f.Name,
		ParentID: f.ParentID,
		UserID:   f.UserID,
	}

	if err := r.db.Save(folderModel).Error; err != nil {
		return err
	}

	f.ID = folderModel.ID
	return nil
}

// FindByID finds a folder by ID
func (r *GormFolderRepository) FindByID(id string) (*folder.Folder, error) {
	var folderModel models.Folder
	if err := r.db.First(&folderModel, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, folder.ErrFolderNotFound
		}
		return nil, err
	}

	return &folder.Folder{
		ID:        folderModel.ID,
		Name:      folderModel.Name,
		ParentID:  folderModel.ParentID,
		UserID:    folderModel.UserID,
		CreatedAt: folderModel.CreatedAt,
		UpdatedAt: folderModel.UpdatedAt,
	}, nil
}

// FindByUserID finds folders by user ID
func (r *GormFolderRepository) FindByUserID(userID string) ([]*folder.Folder, error) {
	var folderModels []models.Folder
	if err := r.db.Where("user_id = ?", userID).Find(&folderModels).Error; err != nil {
		return nil, err
	}

	folders := make([]*folder.Folder, len(folderModels))
	for i, model := range folderModels {
		folders[i] = &folder.Folder{
			ID:        model.ID,
			Name:      model.Name,
			ParentID:  model.ParentID,
			UserID:    model.UserID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		}
	}

	return folders, nil
}

// Delete deletes a folder
func (r *GormFolderRepository) Delete(id string) error {
	return r.db.Delete(&models.Folder{}, "id = ?", id).Error
}

// BelongsToUser checks if a folder belongs to a user
func (r *GormFolderRepository) BelongsToUser(folderID string, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Folder{}).
		Where("id = ? AND user_id = ?", folderID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByUserAndParent finds folders by user ID and parent ID
// If parentID is empty, it returns root-level folders (where ParentID is empty)
func (r *GormFolderRepository) FindByUserAndParent(userID string, parentID string) ([]folder.Folder, error) {
	var folderModels []models.Folder
	query := r.db.Where("user_id = ?", userID)

	if parentID == "" {
		// Find root folders (where parent_id is empty)
		query = query.Where("parent_id IS NULL")
	} else {
		// Find folders with the specified parent
		query = query.Where("parent_id = ?", parentID)
	}

	if err := query.Find(&folderModels).Error; err != nil {
		return nil, err
	}

	folders := make([]folder.Folder, len(folderModels))
	for i, model := range folderModels {
		folders[i] = folder.Folder{
			ID:        model.ID,
			Name:      model.Name,
			ParentID:  model.ParentID,
			UserID:    model.UserID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		}
	}

	return folders, nil
}
