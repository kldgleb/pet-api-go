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
		NewErrorResponse(c, http.StatusUnauthorized, "Empty Authorization header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		NewErrorResponse(c, http.StatusUnauthorized, "Invalid auth header")
		return
	}

	if headerParts[1] == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "Invalid auth header")
		return
	}

	if headerParts[0] != "Bearer" {
		NewErrorResponse(c, http.StatusUnauthorized, "Invalid auth header")
		return
	}

	//parse token
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "Parse token error")
		return
	}
	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	userId, ok := c.Get("userId")
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "error getting user from token")
		return 0, errors.New("user id not found")
	}
	idInt, ok := userId.(int)
	if !ok {
		return 0, errors.New("user id not found")
	}
	return idInt, nil
}
