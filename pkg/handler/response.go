package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	StatusCode int    `json:"StatusCode"`
	Message    string `json:"error"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(statusCode, ":", message)
	c.AbortWithStatusJSON(statusCode, errorResponse{StatusCode: statusCode, Message: message})
}
