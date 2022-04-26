package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/quantum0hound/gochat/internal/handler/server/ws"
	"github.com/quantum0hound/gochat/internal/service"
	"net/http"
)

type Handler struct {
	srv      *service.Service
	wsServer *ws.WebSocketServer
}

func NewHandler(srv *service.Service, wsServer *ws.WebSocketServer) *Handler {
	return &Handler{srv: srv, wsServer: wsServer}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK,
			map[string]interface{}{
				"status": "ok",
			})
	})

	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.signUp)
		auth.POST("/signin", h.signIn)
		auth.POST("/refresh", h.refresh)
	}
	api := router.Group("/api")
	{
		channels := api.Group("/channel", h.userIdentity)
		{
			channels.GET("", h.getAllChannels)
			channels.POST("", h.createChannel)
			channels.GET("/:id/join", h.joinChannel)
			channels.POST("search", h.searchForChannels)
			channels.DELETE("/:id", h.deleteChannel)
			channels.GET("/:id/leave", h.leaveChannel)
		}
		//pubChannels := api.Group("/channel")
		//{
		//	pubChannels.GET("/:id/join", h.joinChannel)
		//}
	}

	return router
}
