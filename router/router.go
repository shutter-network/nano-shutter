package router

import (
	"context"
	"nano-shutter/internal/middleware"
	"nano-shutter/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(ctx context.Context, srv service.Service) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler())

	router.POST("/encrypt", srv.Encrypt)
	router.POST("/decrypt", srv.Decrypt)

	return router
}
