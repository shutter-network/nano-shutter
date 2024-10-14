package router

import (
	"context"
	"nano-shutter/internal/middleware"
	"nano-shutter/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(ctx context.Context, srv service.Service) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler())
	router.Use(cors.Default())

	router.GET("/", srv.GetHandler)
	encrypt := router.Group("/encrypt")
	{
		encrypt.POST("/with_time", srv.EncryptWithTime)
		encrypt.POST("/custom", srv.EncryptCustom)

		// // Creating Options request to disable cors
		// encrypt.OPTIONS("/with_time", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{}) })
		// encrypt.OPTIONS("/custom", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{}) })
	}

	decrypt := router.Group("/decrypt")
	{
		decrypt.POST("/with_time", srv.DecryptWithTime)
		decrypt.POST("/custom", srv.DecryptCustom)

		// // Creating Options request to disable cors
		// encrypt.OPTIONS("/with_time", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{}) })
		// encrypt.OPTIONS("/custom", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{}) })
	}

	return router
}
