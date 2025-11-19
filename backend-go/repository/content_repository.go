package repository

import (
	"backend-go/models"
	"errors"

	"gorm.io/gorm"
)

// ContentRepository define a interface para operações de conteúdo
type ContentRepository interface {
	Create(content *models.Content) error
	GetByID(id uint) (*models.Content, error)
	GetAll(limit, offset int, contentType *string) ([]models.Content, int64, error)
	Update(content *models.Content) error
	Delete(id uint) error
}

type contentRepository struct {
	db *gorm.DB
}

// NewContentRepository cria uma nova instância do ContentRepository
func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{db: db}
}

// Create cria um novo conteúdo no banco de dados
func (r *contentRepository) Create(content *models.Content) error {
	return r.db.Create(content).Error
}

// GetByID busca um conteúdo pelo ID com seus relacionamentos
func (r *contentRepository) GetByID(id uint) (*models.Content, error) {
	var content models.Content
	if err := r.db.Preload("Categories").
		Preload("Interactions").
		First(&content, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("conteúdo não encontrado")
		}
		return nil, err
	}
	return &content, nil
}

// GetAll busca todos os conteúdos com paginação e filtro opcional por tipo
// Retorna a lista de conteúdos e o total de registros
func (r *contentRepository) GetAll(limit, offset int, contentType *string) ([]models.Content, int64, error) {
	var contents []models.Content
	var total int64

	query := r.db.Model(&models.Content{})

	// Aplica filtro por tipo se fornecido
	if contentType != nil && *contentType != "" {
		query = query.Where("type = ?", *contentType)
	}

	// Conta o total de registros (com filtro aplicado)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Busca os registros com paginação
	if err := query.
		Preload("Categories").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&contents).Error; err != nil {
		return nil, 0, err
	}

	return contents, total, nil
}

// Update atualiza um conteúdo existente
func (r *contentRepository) Update(content *models.Content) error {
	return r.db.Save(content).Error
}

// Delete remove um conteúdo do banco de dados
func (r *contentRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Content{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("conteúdo não encontrado")
	}
	return nil
}
