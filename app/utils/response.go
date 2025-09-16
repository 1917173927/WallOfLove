package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type SuccessResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func JsonResponse(c *gin.Context, httpStatusCode int, code int, msg string, data interface{}) {
	c.JSON(httpStatusCode, gin.H{
		"Code": code,
		"Msg":  msg,
		"Data": data,
	})
}

func JsonSuccessResponse(c *gin.Context, data any) {
	JsonResponse(c, http.StatusOK, 200, "OK", data)
}

func JsonErrorResponse(c *gin.Context, code int, msg string) {
	JsonResponse(c, http.StatusOK, 200000 + code, msg, nil)
}

func JsonInternalServerErrorResponse(c *gin.Context) {
	JsonResponse(c, http.StatusInternalServerError, 200500, "Internal server error", nil)
}

