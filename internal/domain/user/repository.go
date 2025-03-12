package user

// Repository defines the interface for user data access
type Repository interface {
	Save(user *User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	UpdateStorageUsed(userID string, storageUsed int64) error
	IncrementStorageUsed(userID string, size int64) error
	DecrementStorageUsed(userID string, size int64) error
}
