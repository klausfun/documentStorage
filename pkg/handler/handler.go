package handler

import (
	"documentStorage/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.NoRoute(func(c *gin.Context) {
		newErrResponse(c, http.StatusNotImplemented, "Method is not implemented")
		return
	})

	api := router.Group("/api")
	{
		api.POST("/register", h.signUp)

		auth := api.Group("/auth")
		{
			auth.POST("/", h.signIn)
			auth.DELETE("/:token", h.signOut)
		}

		docs := api.Group("/docs")
		{
			docs.POST("/", h.createDocument)
			docs.GET("/", h.getAllDocuments)
			docs.GET("/:id", h.getDocumentById)
			docs.DELETE("/:id", h.deleteDocument)
		}
	}

	return router
}
