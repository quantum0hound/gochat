package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/quantum0hound/gochat/internal/models"
	"net/http"
	"time"
)

const (
	refreshTokenCookieName = "RefreshToken"
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

type refreshInput struct {
	Fingerprint string `json:"fingerprint"`
}
type signInInput struct {
	models.User
	refreshInput
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.srv.GetUser(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, err := h.srv.Auth.GenerateAccessToken(user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	session, err := h.srv.CreateSession(user.Id, input.Fingerprint)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	maxAge := session.ExpiresIn.Sub(time.Now())
	c.SetCookie(refreshTokenCookieName, session.RefreshToken, int(maxAge.Seconds()), "", "", false, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": accessToken,
	})

}

func (h *Handler) refresh(c *gin.Context) {
	var input refreshInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//get fingerprint body and refresh token from cookie
	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "refresh token is not provided")
		return
	}

	session, err := h.srv.RefreshSession(refreshToken, input.Fingerprint)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, err := h.srv.Auth.GenerateAccessTokenId(session.UserId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	maxAge := session.ExpiresIn.Sub(time.Now())
	c.SetCookie("RefreshToken", session.RefreshToken, int(maxAge.Seconds()), "", "", false, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": accessToken,
	})

}
