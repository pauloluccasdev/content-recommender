package handler

import (
	"net/http"
	"strconv"

	"backend-go/models"
	"backend-go/service"

	"github.com/gin-gonic/gin"
)

type InteractionHandler struct {
	service              service.InteractionService
	recommendationService service.RecommendationService
}

func NewInteractionHandler(
	service service.InteractionService,
	recommendationService service.RecommendationService,
) *InteractionHandler {
	return &InteractionHandler{
		service:              service,
		recommendationService: recommendationService,
	}
}

func (h *InteractionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("", h.CreateInteraction)
	rg.GET("/user/:user_id", h.GetUserInteractions)
	rg.GET("/content/:content_id", h.GetContentInteractions)
}

// DTOs de Request
type CreateInteractionRequest struct {
	UserID          uint     `json:"user_id" binding:"required"`
	ContentID       uint     `json:"content_id" binding:"required"`
	InteractionType string   `json:"interaction_type" binding:"required,oneof=view like dislike rating share comment"`
	Rating          *float64 `json:"rating,omitempty"`
}

// DTO de Response
type InteractionResponse struct {
	ID              uint     `json:"id"`
	UserID          uint     `json:"user_id"`
	ContentID       uint     `json:"content_id"`
	InteractionType string   `json:"interaction_type"`
	Rating          *float64 `json:"rating,omitempty"`
	CreatedAt       string   `json:"created_at"`
}

func newInteractionResponse(interaction *models.UserInteraction) InteractionResponse {
	return InteractionResponse{
		ID:              interaction.ID,
		UserID:          interaction.UserID,
		ContentID:       interaction.ContentID,
		InteractionType: interaction.InteractionType,
		Rating:          interaction.Rating,
		CreatedAt:       interaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CreateInteraction godoc
// @Summary Cria uma nova interação de usuário com conteúdo
// @Tags interactions
// @Accept json
// @Produce json
// @Param interaction body CreateInteractionRequest true "Dados da interação"
// @Success 201 {object} InteractionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /interactions [post]
func (h *InteractionHandler) CreateInteraction(c *gin.Context) {
	var req CreateInteractionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	interaction, err := h.service.CreateInteraction(
		req.UserID,
		req.ContentID,
		req.InteractionType,
		req.Rating,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newInteractionResponse(interaction))
	
	// Notificar motor Python em background (não bloqueia resposta)
	go func() {
		_ = h.recommendationService.NotifyNewInteraction(
			req.UserID,
			req.ContentID,
			req.InteractionType,
			req.Rating,
		)
	}()
}

// GetUserInteractions godoc
// @Summary Lista todas as interações de um usuário
// @Tags interactions
// @Produce json
// @Param user_id path int true "ID do usuário"
// @Success 200 {array} InteractionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /interactions/user/{user_id} [get]
func (h *InteractionHandler) GetUserInteractions(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	interactions, err := h.service.GetUserInteractions(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]InteractionResponse, len(interactions))
	for i := range interactions {
		responses[i] = newInteractionResponse(&interactions[i])
	}

	c.JSON(http.StatusOK, responses)
}

// GetContentInteractions godoc
// @Summary Lista todas as interações de um conteúdo
// @Tags interactions
// @Produce json
// @Param content_id path int true "ID do conteúdo"
// @Success 200 {array} InteractionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /interactions/content/{content_id} [get]
func (h *InteractionHandler) GetContentInteractions(c *gin.Context) {
	contentIDParam := c.Param("content_id")
	contentID, err := strconv.ParseUint(contentIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de conteúdo inválido"})
		return
	}

	interactions, err := h.service.GetContentInteractions(uint(contentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]InteractionResponse, len(interactions))
	for i := range interactions {
		responses[i] = newInteractionResponse(&interactions[i])
	}

	c.JSON(http.StatusOK, responses)
}

