package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type error struct {
	StatusCode int    `json:"StatusCode"`
	Message    string `json:"Message"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(statusCode, ": ", message)
	c.AbortWithStatusJSON(statusCode, error{StatusCode: statusCode, Message: message})
}
