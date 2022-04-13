package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/quantum0hound/gochat/internal/models"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.srv.CreateUser(&user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {

}
