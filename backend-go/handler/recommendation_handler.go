package handler

import (
	"net/http"
	"strconv"

	"backend-go/service"

	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	service service.RecommendationService
}

func NewRecommendationHandler(service service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{service: service}
}

func (h *RecommendationHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/user/:user_id", h.GetRecommendations)
}

// DTO de Request
type GetRecommendationsRequest struct {
	TopN   int    `form:"top_n" binding:"omitempty,min=1,max=50"`
	Method string `form:"method" binding:"omitempty,oneof=similarity popularity"`
}

// DTO de Response
type RecommendationResponse struct {
	UserID         uint   `json:"user_id"`
	ContentIDs     []uint `json:"content_ids"`
	Method         string `json:"method"`
	Count          int    `json:"count"`
}

// GetRecommendations godoc
// @Summary Obtém recomendações para um usuário
// @Tags recommendations
// @Produce json
// @Param user_id path int true "ID do usuário"
// @Param top_n query int false "Número de recomendações" default(10) minimum(1) maximum(50)
// @Param method query string false "Método de recomendação" Enums(similarity, popularity) default(similarity)
// @Success 200 {object} RecommendationResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /recommendations/user/{user_id} [get]
func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	var req GetRecommendationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Valores padrão
	topN := req.TopN
	if topN == 0 {
		topN = 10
	}
	method := req.Method
	if method == "" {
		method = "similarity"
	}

	contentIDs, err := h.service.GetRecommendations(uint(userID), topN, method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao obter recomendações",
			"details": err.Error(),
		})
		return
	}

	response := RecommendationResponse{
		UserID:     uint(userID),
		ContentIDs: contentIDs,
		Method:     method,
		Count:      len(contentIDs),
	}

	c.JSON(http.StatusOK, response)
}

