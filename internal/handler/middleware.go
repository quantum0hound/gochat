package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headersPart := strings.Split(header, " ")
	if len(headersPart) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "incorrect auth header")
		return
	}

	userId, err := h.srv.Auth.ParseAccessToken(headersPart[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, userId)

}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		msg := "user id is not found"
		newErrorResponse(c, http.StatusInternalServerError, msg)
		return 0, errors.New(msg)
	}
	idInt, ok := id.(int)
	if !ok {
		msg := "user id is of invalid type"
		newErrorResponse(c, http.StatusInternalServerError, msg)
		return 0, errors.New(msg)
	}

	return idInt, nil
}
