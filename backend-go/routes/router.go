package routes

import (
	"backend-go/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	engine         *gin.Engine
	userHandler    *handler.UserHandler
	contentHandler *handler.ContentHandler
}

func NewRouter(userHandler *handler.UserHandler, contentHandler *handler.ContentHandler) *Router {
	engine := gin.Default()
	return &Router{
		engine:         engine,
		userHandler:    userHandler,
		contentHandler: contentHandler,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.engine.Group("/api")
	
	// Rotas de usuários
	users := api.Group("/users")
	r.userHandler.RegisterRoutes(users)
	
	// Rotas de conteúdos
	contents := api.Group("/contents")
	r.contentHandler.RegisterRoutes(contents)
	
	return r.engine
}
