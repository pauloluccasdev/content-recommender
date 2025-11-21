package service

import (
	"backend-go/models"
	"backend-go/repository"
	"errors"
	"time"
)

// InteractionService define a interface para operações de interações
type InteractionService interface {
	CreateInteraction(userID, contentID uint, interactionType string, rating *float64) (*models.UserInteraction, error)
	GetUserInteractions(userID uint) ([]models.UserInteraction, error)
	GetContentInteractions(contentID uint) ([]models.UserInteraction, error)
}

type interactionService struct {
	repo repository.InteractionRepository
}

// NewInteractionService cria uma nova instância do InteractionService
func NewInteractionService(repo repository.InteractionRepository) InteractionService {
	return &interactionService{repo: repo}
}

// CreateInteraction cria uma nova interação
func (s *interactionService) CreateInteraction(userID, contentID uint, interactionType string, rating *float64) (*models.UserInteraction, error) {
	// Validar tipo de interação
	validTypes := map[string]bool{
		"view":     true,
		"like":     true,
		"dislike":  true,
		"rating":   true,
		"share":    true,
		"comment":  true,
	}
	if !validTypes[interactionType] {
		return nil, errors.New("tipo de interação inválido")
	}

	// Validar rating se fornecido
	if rating != nil {
		if *rating < 1.0 || *rating > 5.0 {
			return nil, errors.New("rating deve estar entre 1.0 e 5.0")
		}
		// Se o tipo não for rating, mas rating foi fornecido, usar tipo rating
		if interactionType != "rating" {
			interactionType = "rating"
		}
	}

	// Se for tipo rating mas não tem rating, retornar erro
	if interactionType == "rating" && rating == nil {
		return nil, errors.New("rating é obrigatório para interação do tipo rating")
	}

	interaction := &models.UserInteraction{
		UserID:          userID,
		ContentID:       contentID,
		InteractionType: interactionType,
		Rating:          rating,
		CreatedAt:       time.Now(),
	}

	if err := s.repo.Create(interaction); err != nil {
		return nil, err
	}

	return interaction, nil
}

// GetUserInteractions retorna todas as interações de um usuário
func (s *interactionService) GetUserInteractions(userID uint) ([]models.UserInteraction, error) {
	return s.repo.GetByUserID(userID)
}

// GetContentInteractions retorna todas as interações de um conteúdo
func (s *interactionService) GetContentInteractions(contentID uint) ([]models.UserInteraction, error) {
	return s.repo.GetByContentID(contentID)
}

