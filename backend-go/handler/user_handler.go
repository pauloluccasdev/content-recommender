package handler

import (
	"net/http"
	"time"

	"backend-go/models"
	"backend-go/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.GetFirstUser)
}

type userResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(u *models.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

// GetFirstUser godoc
// @Summary Retorna o primeiro usuário seed
// @Tags users
// @Produce json
// @Success 200 {object} userResponse
// @Failure 404 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) GetFirstUser(c *gin.Context) {
	user, err := h.service.GetFirstUser()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newUserResponse(user))
}
