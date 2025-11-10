package routes

import (
	"backend-go/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	engine      *gin.Engine
	userHandler *handler.UserHandler
}

func NewRouter(userHandler *handler.UserHandler) *Router {
	engine := gin.Default()
	return &Router{engine: engine, userHandler: userHandler}
}

func (r *Router) SetupRoutes() *gin.Engine {
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.engine.Group("/api")
	users := api.Group("/users")
	r.userHandler.RegisterRoutes(users)
	return r.engine
}
