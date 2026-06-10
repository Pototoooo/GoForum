package response

import (
	"net/http"

	"bluebell/pkg/Code"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code Code.ResponseCode `json:"code"`
	Msg  string            `json:"msg"`
	Data interface{}       `json:"data,omitempty"`
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: Code.CodeSuccess,
		Msg:  Code.CodeSuccess.Msg(),
		Data: data,
	})
}

func ResponseError(c *gin.Context, code Code.ResponseCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

// ResponseErrorWithMsg 返回带有自定义错误消息的错误响应
func ResponseErrorWithMsg(c *gin.Context, code Code.ResponseCode, msg string) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
