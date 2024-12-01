package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func (e *ErrorResponse) Error() string {
	return e.Text
}

func NewErrorResponse(code int, text string) *ErrorResponse {
	return &ErrorResponse{Code: code, Text: text}
}

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
	errResp := ErrorResponse{statusCode, message}
	c.AbortWithStatusJSON(statusCode, errorModel{errResp})
}
