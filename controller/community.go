package controller

import (
	"strconv"

	"GoForum/logic"
	"GoForum/pkg/Code"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 查询社区列表
// @Summary     社区列表
// @Description 获取所有社区列表
// @Tags       社区
// @Accept     json
// @Produce    json
// @Success   200 {object} response.ResponseData{data=[]models.Community}
// @Router    /community [get]
func CommunityHandler(c *gin.Context) {
	data, err := logic.GetCommunity()
	if err != nil {
		zap.L().Error("logic.CommunityHandler() error", zap.Error(err))
		ResponseError(c, Code.CodeServiceUnavailable)
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 查询社区详情
// @Summary     社区详情
// @Description 根据社区 ID 获取社区详细信息
// @Tags       社区
// @Accept     json
// @Produce    json
// @Param      id path int true "社区 ID"
// @Success   200 {object} response.ResponseData{data=models.DetailCommunity}
// @Router    /community/{id} [get]
func CommunityDetailHandler(c *gin.Context) {
	communityID := c.Param("id")
	id, err := strconv.ParseInt(communityID, 10, 64)
	if err != nil {
		ResponseError(c, Code.CodeInvalidParam)
		return
	}
	data, err := logic.GetDetailCommunity(id)
	if err != nil {
		zap.L().Error("logic.CommunityDetailHandler() error", zap.Error(err))
		ResponseError(c, Code.CodeServiceUnavailable)
		return
	}
	ResponseSuccess(c, data)
}
