package handler

import (
	"documentStorage/pkg"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type responseModel struct {
	Response any `json:"response"`
}

type dataModel struct {
	Data any `json:"data"`
}

type errorModel struct {
	Err any `json:"error"`
}

func newErrResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorModel{
		pkg.ErrorResponse{
			Code: statusCode,
			Text: message,
		},
	})
}
