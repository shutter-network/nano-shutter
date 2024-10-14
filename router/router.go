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

	encrypt := router.Group("/encrypt")
	{
		encrypt.POST("/with_time", srv.EncryptWithTime)
		encrypt.POST("/custom", srv.EncryptCustom)
	}

	decrypt := router.Group("/decrypt")
	{
		decrypt.POST("/with_time", srv.DecryptWithTime)
		decrypt.POST("/custom", srv.DecryptCustom)
	}

	return router
}
