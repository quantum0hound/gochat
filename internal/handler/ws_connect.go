package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) connectWebSocket(c *gin.Context) {
	h.wsServer.ServePeer(c.Writer, c.Request)
}
