package handler

import (
	"net/http"
	"strconv"
	"time"

	"backend-go/models"
	"backend-go/service"

	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	service service.ContentService
}

func NewContentHandler(service service.ContentService) *ContentHandler {
	return &ContentHandler{service: service}
}

// RegisterRoutes registra as rotas do handler de conteúdo
func (h *ContentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("", h.CreateContent)
	rg.GET("", h.ListContents)
	rg.GET("/:id", h.GetContentByID)
	rg.PUT("/:id", h.UpdateContent)
	rg.DELETE("/:id", h.DeleteContent)
}

// DTOs de Request
type CreateContentRequest struct {
	Title       string    `json:"title" binding:"required,min=3,max=200"`
	Description string    `json:"description" binding:"max=1000"`
	Type        string    `json:"type" binding:"required,oneof=article video podcast book course"`
	ReleaseDate time.Time `json:"release_date" binding:"required"`
	CategoryIDs []uint    `json:"category_ids,omitempty"`
}

type UpdateContentRequest struct {
	Title       *string    `json:"title,omitempty" binding:"omitempty,min=3,max=200"`
	Description *string    `json:"description,omitempty" binding:"omitempty,max=1000"`
	Type        *string    `json:"type,omitempty" binding:"omitempty,oneof=article video podcast book course"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	CategoryIDs []uint     `json:"category_ids,omitempty"`
}

// DTO de Response
type ContentResponse struct {
	ID          uint               `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Type        string             `json:"type"`
	ReleaseDate time.Time          `json:"release_date"`
	CreatedAt   time.Time          `json:"created_at"`
	Categories  []CategoryResponse `json:"categories,omitempty"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ListContentsResponse struct {
	Contents []ContentResponse `json:"contents"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

// Helper para converter model em response DTO
func newContentResponse(c *models.Content) ContentResponse {
	categories := make([]CategoryResponse, len(c.Categories))
	for i, cat := range c.Categories {
		categories[i] = CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
		}
	}

	return ContentResponse{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
		Type:        c.Type,
		ReleaseDate: c.ReleaseDate,
		CreatedAt:   c.CreatedAt,
		Categories:  categories,
	}
}

// CreateContent godoc
// @Summary Cria um novo conteúdo
// @Tags contents
// @Accept json
// @Produce json
// @Param content body CreateContentRequest true "Dados do conteúdo"
// @Success 201 {object} ContentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /contents [post]
func (h *ContentHandler) CreateContent(c *gin.Context) {
	var req CreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converte para o request do service
	serviceReq := service.CreateContentRequest{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		ReleaseDate: req.ReleaseDate,
		CategoryIDs: req.CategoryIDs,
	}

	content, err := h.service.CreateContent(serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newContentResponse(content))
}

// GetContentByID godoc
// @Summary Busca um conteúdo pelo ID
// @Tags contents
// @Produce json
// @Param id path int true "ID do conteúdo"
// @Success 200 {object} ContentResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /contents/{id} [get]
func (h *ContentHandler) GetContentByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	content, err := h.service.GetContentByID(uint(id))
	if err != nil {
		if err.Error() == "conteúdo não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newContentResponse(content))
}

// ListContents godoc
// @Summary Lista conteúdos com paginação e filtros
// @Tags contents
// @Produce json
// @Param page query int false "Número da página" default(1)
// @Param limit query int false "Itens por página" default(20)
// @Param type query string false "Filtro por tipo (article, video, podcast, book, course)"
// @Success 200 {object} ListContentsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /contents [get]
func (h *ContentHandler) ListContents(c *gin.Context) {
	// Parse dos query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var contentType *string
	if typeParam := c.Query("type"); typeParam != "" {
		contentType = &typeParam
	}

	contents, total, err := h.service.ListContents(page, limit, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Converte para response DTOs
	contentResponses := make([]ContentResponse, len(contents))
	for i, content := range contents {
		contentResponses[i] = newContentResponse(&content)
	}

	response := ListContentsResponse{
		Contents: contentResponses,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateContent godoc
// @Summary Atualiza um conteúdo existente
// @Tags contents
// @Accept json
// @Produce json
// @Param id path int true "ID do conteúdo"
// @Param content body UpdateContentRequest true "Dados para atualização"
// @Success 200 {object} ContentResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /contents/{id} [put]
func (h *ContentHandler) UpdateContent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converte para o request do service
	serviceReq := service.UpdateContentRequest{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		ReleaseDate: req.ReleaseDate,
		CategoryIDs: req.CategoryIDs,
	}

	content, err := h.service.UpdateContent(uint(id), serviceReq)
	if err != nil {
		if err.Error() == "conteúdo não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newContentResponse(content))
}

// DeleteContent godoc
// @Summary Remove um conteúdo
// @Tags contents
// @Produce json
// @Param id path int true "ID do conteúdo"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /contents/{id} [delete]
func (h *ContentHandler) DeleteContent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.service.DeleteContent(uint(id)); err != nil {
		if err.Error() == "conteúdo não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "conteúdo removido com sucesso"})
}
