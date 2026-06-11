package controller

import (
	"GoForum/logic"
	"GoForum/pkg/Code"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Vote struct {
	PostID    int64 `json:"post_id,string" binding:"required"`
	Direction int   `json:"direction" binding:"required,oneof=0 1 -1"`
}

// VoteHandler 投票
// @Summary     帖子投票
// @Description 对帖子进行投票（1-赞同, -1-反对, 0-取消投票），需要登录
// @Tags       投票
// @Accept     json
// @Produce    json
// @Security   BearerAuth
// @Param      vote body Vote true "投票参数"
// @Success   200 {object} response.ResponseData
// @Router    /vote [post]
func VoteHandler(c *gin.Context) {
	// 参数校验
	vote := new(Vote)
	err := c.ShouldBindJSON(vote)
	if err != nil {
		ResponseError(c, Code.CodeInvalidParam)
		zap.L().Error("VoteHandler ShouldBindJSON failed", zap.Error(err))
		return
	}
	// 拿到投票用户
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, Code.CodeNeedLogin)
		zap.L().Error("VoteHandler GetCurrentUserID failed", zap.Error(err))
		return
	}
	// 调用业务逻辑
	err = logic.VoteHandler(userID, vote.Direction, vote.PostID)
	if err != nil {
		ResponseError(c, Code.CodeServiceUnavailable)
		zap.L().Error("VoteHandler logic.VoteHandler failed", zap.Error(err))
		return
	}
	ResponseSuccess(c, nil)
}
