package service

import (
	"errors"
	"backend-go/models"
	"backend-go/repository"
	"time"
)

// ContentService define a interface para operações de negócio de conteúdo
type ContentService interface {
	CreateContent(req CreateContentRequest) (*models.Content, error)
	GetContentByID(id uint) (*models.Content, error)
	ListContents(page, limit int, contentType *string) ([]models.Content, int64, error)
	UpdateContent(id uint, req UpdateContentRequest) (*models.Content, error)
	DeleteContent(id uint) error
}

type contentService struct {
	repo repository.ContentRepository
}

// CreateContentRequest representa os dados necessários para criar um conteúdo
type CreateContentRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	ReleaseDate time.Time `json:"release_date"`
	CategoryIDs []uint    `json:"category_ids,omitempty"`
}

// UpdateContentRequest representa os dados necessários para atualizar um conteúdo
type UpdateContentRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Type        *string    `json:"type,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	CategoryIDs []uint     `json:"category_ids,omitempty"`
}

// NewContentService cria uma nova instância do ContentService
func NewContentService(repo repository.ContentRepository) ContentService {
	return &contentService{repo: repo}
}

// Tipos de conteúdo válidos
var validContentTypes = map[string]bool{
	"article":  true,
	"video":    true,
	"podcast":  true,
	"book":     true,
	"course":   true,
}

// CreateContent cria um novo conteúdo após validar os dados
func (s *contentService) CreateContent(req CreateContentRequest) (*models.Content, error) {
	// Validações de negócio
	if err := s.validateContentRequest(req.Title, req.Description, req.Type, req.ReleaseDate); err != nil {
		return nil, err
	}

	content := &models.Content{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		ReleaseDate: req.ReleaseDate,
	}

	// Se houver categorias, associa (será feito via relacionamento many-to-many)
	// Por enquanto, criamos o conteúdo primeiro
	if err := s.repo.Create(content); err != nil {
		return nil, err
	}

	return content, nil
}

// GetContentByID busca um conteúdo pelo ID
func (s *contentService) GetContentByID(id uint) (*models.Content, error) {
	if id == 0 {
		return nil, errors.New("ID inválido")
	}

	content, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// ListContents lista conteúdos com paginação e filtros
func (s *contentService) ListContents(page, limit int, contentType *string) ([]models.Content, int64, error) {
	// Valores padrão e validações
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100 // Limite máximo
	}

	offset := (page - 1) * limit

	contents, total, err := s.repo.GetAll(limit, offset, contentType)
	if err != nil {
		return nil, 0, err
	}

	return contents, total, nil
}

// UpdateContent atualiza um conteúdo existente
func (s *contentService) UpdateContent(id uint, req UpdateContentRequest) (*models.Content, error) {
	if id == 0 {
		return nil, errors.New("ID inválido")
	}

	// Busca o conteúdo existente
	content, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Atualiza apenas os campos fornecidos
	if req.Title != nil {
		if *req.Title == "" {
			return nil, errors.New("título não pode ser vazio")
		}
		if len(*req.Title) > 200 {
			return nil, errors.New("título deve ter no máximo 200 caracteres")
		}
		content.Title = *req.Title
	}

	if req.Description != nil {
		if len(*req.Description) > 1000 {
			return nil, errors.New("descrição deve ter no máximo 1000 caracteres")
		}
		content.Description = *req.Description
	}

	if req.Type != nil {
		if !validContentTypes[*req.Type] {
			return nil, errors.New("tipo de conteúdo inválido")
		}
		content.Type = *req.Type
	}

	if req.ReleaseDate != nil {
		content.ReleaseDate = *req.ReleaseDate
	}

	// Atualiza no banco
	if err := s.repo.Update(content); err != nil {
		return nil, err
	}

	return content, nil
}

// DeleteContent remove um conteúdo
func (s *contentService) DeleteContent(id uint) error {
	if id == 0 {
		return errors.New("ID inválido")
	}

	return s.repo.Delete(id)
}

// validateContentRequest valida os dados de criação de conteúdo
func (s *contentService) validateContentRequest(title, description, contentType string, releaseDate time.Time) error {
	if title == "" {
		return errors.New("título é obrigatório")
	}
	if len(title) < 3 {
		return errors.New("título deve ter no mínimo 3 caracteres")
	}
	if len(title) > 200 {
		return errors.New("título deve ter no máximo 200 caracteres")
	}

	if len(description) > 1000 {
		return errors.New("descrição deve ter no máximo 1000 caracteres")
	}

	if contentType == "" {
		return errors.New("tipo de conteúdo é obrigatório")
	}
	if !validContentTypes[contentType] {
		return errors.New("tipo de conteúdo inválido. Tipos válidos: article, video, podcast, book, course")
	}

	// Valida que a data não seja muito no futuro (opcional, mas boa prática)
	if releaseDate.After(time.Now().AddDate(10, 0, 0)) {
		return errors.New("data de lançamento não pode ser mais de 10 anos no futuro")
	}

	return nil
}

