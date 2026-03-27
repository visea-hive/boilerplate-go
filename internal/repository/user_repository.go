package repository

import (
	"github.com/visea-hive/auth-core/internal/model"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user-related database operations.
type UserRepository interface {
	Create(db *gorm.DB, user *model.User) error
	GetByID(db *gorm.DB, id uint) (*model.User, error)
	GetByUUID(db *gorm.DB, uuid string) (*model.User, error)
	GetByEmail(db *gorm.DB, email string) (*model.User, error)
	Update(db *gorm.DB, user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(db *gorm.DB, user *model.User) error {
	return db.Create(user).Error
}

func (r *userRepository) GetByID(db *gorm.DB, id uint) (*model.User, error) {
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUUID(db *gorm.DB, uuid string) (*model.User, error) {
	var user model.User
	if err := db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(db *gorm.DB, email string) (*model.User, error) {
	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(db *gorm.DB, user *model.User) error {
	return db.Save(user).Error
}
