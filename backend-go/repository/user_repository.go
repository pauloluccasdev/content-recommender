package repository

import (
	"backend-go/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	First() (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) First() (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Recommendations").
		Preload("Interactions").
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
