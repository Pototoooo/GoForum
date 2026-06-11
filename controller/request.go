package controller

import (
	"GoForum/middlewire"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ErrorRounterNoUserInfo = errors.New("路由无用户信息")

func GetCurrentUserID(c *gin.Context) (userID int64, err error) {
	userIDValue, ok := c.Get(middlewire.UserID)
	if !ok {
		err = ErrorRounterNoUserInfo
		return
	}
	userID, ok = userIDValue.(int64)
	if !ok {
		err = errors.New("用户信息类型错误")
		return
	}
	return
}

// 拿到页码和一页的size
func GetPageNSize(c *gin.Context) (page, size int) {
	sizeStr := c.Query("size")
	pageStr := c.Query("page")
	size, _ = strconv.Atoi(sizeStr)
	page, _ = strconv.Atoi(pageStr)
	return
}
