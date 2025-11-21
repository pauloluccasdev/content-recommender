package repository

import (
	"backend-go/models"

	"gorm.io/gorm"
)

// InteractionRepository define a interface para operações de interações
type InteractionRepository interface {
	Create(interaction *models.UserInteraction) error
	GetByUserID(userID uint) ([]models.UserInteraction, error)
	GetByContentID(contentID uint) ([]models.UserInteraction, error)
	GetByUserAndContent(userID, contentID uint) (*models.UserInteraction, error)
	CountByUserID(userID uint) (int64, error)
}

type interactionRepository struct {
	db *gorm.DB
}

// NewInteractionRepository cria uma nova instância do InteractionRepository
func NewInteractionRepository(db *gorm.DB) InteractionRepository {
	return &interactionRepository{db: db}
}

// Create cria uma nova interação no banco de dados
func (r *interactionRepository) Create(interaction *models.UserInteraction) error {
	return r.db.Create(interaction).Error
}

// GetByUserID busca todas as interações de um usuário
func (r *interactionRepository) GetByUserID(userID uint) ([]models.UserInteraction, error) {
	var interactions []models.UserInteraction
	if err := r.db.Where("user_id = ?", userID).
		Preload("Content").
		Order("created_at DESC").
		Find(&interactions).Error; err != nil {
		return nil, err
	}
	return interactions, nil
}

// GetByContentID busca todas as interações de um conteúdo
func (r *interactionRepository) GetByContentID(contentID uint) ([]models.UserInteraction, error) {
	var interactions []models.UserInteraction
	if err := r.db.Where("content_id = ?", contentID).
		Preload("User").
		Order("created_at DESC").
		Find(&interactions).Error; err != nil {
		return nil, err
	}
	return interactions, nil
}

// GetByUserAndContent busca uma interação específica de um usuário com um conteúdo
func (r *interactionRepository) GetByUserAndContent(userID, contentID uint) (*models.UserInteraction, error) {
	var interaction models.UserInteraction
	if err := r.db.Where("user_id = ? AND content_id = ?", userID, contentID).
		Order("created_at DESC").
		First(&interaction).Error; err != nil {
		return nil, err
	}
	return &interaction, nil
}

// CountByUserID conta o número de interações de um usuário
func (r *interactionRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.UserInteraction{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

