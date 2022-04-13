package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/quantum0hound/gochat/internal/service"
	"net/http"
)

type Handler struct {
	srv *service.Service
}

func NewHandler(srv *service.Service) *Handler {
	return &Handler{srv: srv}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}

	router.Use(cors.New(config))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK,
			map[string]interface{}{
				"status": "ok",
			})
	})

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	channels := router.Group("/channel")
	{
		channels.GET("", h.getAllChannels)
		channels.POST("", h.createChannel)
		channels.GET("/:id/join", h.joinChannel)
		channels.GET("/:id/leave", h.leaveChannel)
		channels.DELETE("/:id", h.deleteChannel)
	}

	return router
}
