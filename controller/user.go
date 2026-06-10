package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models/param"
	"bluebell/pkg/Code"
	"bluebell/pkg/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterHandler 注册用户
// @Summary     用户注册
// @Description 注册新用户账号
// @Tags       用户
// @Accept     json
// @Produce    json
// @Param      params body param.RegisterParams true "注册参数"
// @Success   200 {object} response.ResponseData
// @Router    /signup [post]
func RegisterHandler(c *gin.Context) {
	params := new(param.RegisterParams)
	if err := c.ShouldBindJSON(params); err != nil {
		ResponseError(c, Code.CodeInvalidParam)
		zap.L().Error("参数绑定失败", zap.Error(err))
		return
	}
	err := logic.Register(params)
	if err != nil {
		ResponseError(c, Code.CodeInvalidParam)
		zap.L().Error("用户注册失败", zap.Error(err), zap.String("username", params.Username))
		return
	}
	ResponseSuccess(c, nil)
}

// LoginHandler 用户登录
// @Summary     用户登录
// @Description 用户登录，获取 JWT Token
// @Tags       用户
// @Accept     json
// @Produce    json
// @Param      params body param.LoginParams true "登录参数"
// @Success   200 {object} response.ResponseData
// @Router    /login [post]
func LoginHandler(c *gin.Context) {
	params := new(param.LoginParams)
	if err := c.ShouldBindJSON(params); err != nil {
		ResponseError(c, Code.CodeInvalidParam)
		zap.L().Error("参数绑定失败", zap.Error(err))
		return
	}
	err := logic.Login(params)
	if err != nil {
		ResponseError(c, Code.CodeInvalidParam)
		zap.L().Error("用户登录失败", zap.Error(err), zap.String("username", params.Username))
		return
	}
	userID := mysql.QueryUserIdByName(params.Username)
	aToken, _ := jwt.GenToken(userID, params.Username)
	ResponseSuccess(c, gin.H{
		"username": params.Username,
		"token":    aToken,
	})
}
