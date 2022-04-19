package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/quantum0hound/gochat/internal/models"
	"net/http"
)

func (h *Handler) getAllChannels(c *gin.Context) {

}

func (h *Handler) createChannel(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		return
	}

	var channel models.Channel
	if err := c.BindJSON(&channel); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	channel.Creator = id

	channelId, err := h.srv.Create(&channel)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": channelId,
	})
}

func (h *Handler) joinChannel(c *gin.Context) {

}

func (h *Handler) leaveChannel(c *gin.Context) {

}

func (h *Handler) deleteChannel(c *gin.Context) {

}
