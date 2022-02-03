package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test-api/pkg/entity"
)

func (h *Handler) signUp(c *gin.Context) {
	var input entity.User
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.CreateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.services.Authorization.GetJWTByCredentials(input.Username, input.Password)

	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"token": token,
	})
}
