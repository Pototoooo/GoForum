package controller

import (
	"errors"
	"strconv"

	"GoForum/logic"
	"GoForum/models/param"
	"GoForum/models/post"
	"GoForum/pkg/Code"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PostHandler 创建帖子
// @Summary     创建帖子
// @Description 创建一个新帖子（需要登录）
// @Tags       帖子
// @Accept     json
// @Produce    json
// @Security   BearerAuth
// @Param      post body post.Post true "帖子信息"
// @Success   200 {object} response.ResponseData
// @Router    /post [post]
func PostHandler(c *gin.Context) {
	p := new(post.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, Code.CodeInvalidParam)
		zap.L().Error("invalid param", zap.Error(err))
		return
	}
	// 校验社区ID是否有效
	if p.CommunityID <= 0 {
		ResponseError(c, Code.CodeInvalidParam)
		zap.L().Error("invalid community_id", zap.Int64("community_id", p.CommunityID))
		return
	}
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, Code.CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	if err := logic.CreatePost(p); err != nil {
		if errors.Is(err, logic.ErrorCommunityNotFound) {
			ResponseError(c, Code.CodeNotFound)
			zap.L().Error("community not found", zap.Int64("community_id", p.CommunityID))
			return
		}
		zap.L().Error("logic.CreatePost() failed", zap.Error(err))
		ResponseError(c, Code.CodeServiceUnavailable)
		return
	}
	ResponseSuccess(c, nil)
}

// GetDetailPostHandler 查询帖子详情
// @Summary     帖子详情
// @Description 根据帖子 ID 获取帖子详细信息
// @Tags       帖子
// @Accept     json
// @Produce    json
// @Param      id path int true "帖子 ID"
// @Success   200 {object} response.ResponseData{data=post.DetailPost}
// @Router    /post/{id} [get]
func GetDetailPostHandler(c *gin.Context) {
	// 从url获取帖子id
	postID := c.Param("id")
	if postID == "" {
		ResponseError(c, Code.CodeInvalidParam)
		return
	}

	// 将postID解析为int64
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		ResponseError(c, Code.CodeInvalidParam)
		return
	}
	// 数据库由id查数据
	data, err := logic.GetPostByID(id)
	if err != nil {
		zap.L().Error("logic.GetPostByID() failed", zap.Error(err))
		ResponseError(c, Code.CodeServiceUnavailable)
		return
	}
	// 返回
	ResponseSuccess(c, data)
}

// GetPostListHandler 查询帖子列表（分页）
// @Summary     帖子列表（分页）
// @Description 按时间倒序分页获取帖子列表
// @Tags       帖子
// @Accept     json
// @Produce    json
// @Param      page query int false "页码" default(1)
// @Param      size query int false "每页条数" default(10)
// @Success   200 {object} response.ResponseData{data=[]post.DetailPost}
// @Router    /posts [get]
func GetPostListHandler(c *gin.Context) {
	// 拿到页数据
	page, size := GetPageNSize(c)
	// 调用逻辑
	data, err := logic.GetPostList(page, size)
	if err != nil {
		ResponseError(c, Code.CodeServiceUnavailable)
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		return
	}
	ResponseSuccess(c, data)
}

// GetOrderedPosts 查询有序帖子列表
// @Summary     有序帖子列表
// @Description 支持按时间/分数排序，可按社区筛选
// @Tags       帖子
// @Accept     json
// @Produce    json
// @Param      page query int false "页码" default(1)
// @Param      size query int false "每页条数" default(10)
// @Param      order query string false "排序方式：time(时间)/score(分数)" default(time)
// @Param      community_id query int false "社区 ID（可选）"
// @Success   200 {object} response.ResponseData{data=[]post.DetailPost}
// @Router    /posts2 [get]
func GetOrderedPosts(c *gin.Context) {
	param := &param.PostsPageParams{
		Size:  10,
		Page:  1,
		Order: "time",
	}
	c.ShouldBindQuery(param)
	data, err := logic.GetPostListByOrder(param)
	if err != nil {
		if errors.Is(err, logic.ErrorCommunityNotFound) {
			ResponseErrorWithMsg(c, Code.CodeNotFound, "社区不存在")
			return
		}
		zap.L().Error("logic.GetPostListByOrder() failed", zap.Error(err))
		ResponseError(c, Code.CodeServiceUnavailable)
		return
	}
	ResponseSuccess(c, data)
}

// func CommuPostsHandler(c *gin.Context) {
// 	// 初始化默认参数
// 	param := &param.PostsPageParams{
// 		CommunityID: 2213,
// 		Size:        2,
// 		Page:        1,
// 		Order:       "time",
// 	}
// 	// 获取请求头参数
// 	c.ShouldBindQuery(param)
// 	// 调用逻辑
// 	data, err := logic.GetPostListByOrder(param)
// 	if err != nil {
// 		zap.L().Error("logic.GetPostListByOrder() failed", zap.Error(err))
// 		ResponseError(c, Code.CodeServiceUnavailable)
// 		return
// 	}
// 	ResponseSuccess(c, data)
// }
