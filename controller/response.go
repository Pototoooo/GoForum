package controller

import (
	"bluebell/pkg/Code"
	"bluebell/pkg/response"

	"github.com/gin-gonic/gin"
)

func ResponseSuccess(c *gin.Context, data interface{}) {
	response.ResponseSuccess(c, data)
}

func ResponseError(c *gin.Context, code Code.ResponseCode) {
	response.ResponseError(c, code)
}

func ResponseErrorWithMsg(c *gin.Context, code Code.ResponseCode, msg string) {
	response.ResponseErrorWithMsg(c, code, msg)
}
