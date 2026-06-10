// Package logic 封装帖子相关的业务逻辑层。
package logic

import (
	"errors"

	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models/param"
	"bluebell/models/post"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

var ErrorCommunityNotFound = errors.New("社区不存在")

func CreatePost(p *post.Post) (err error) {
	// 校验社区是否存在
	if _, err = mysql.GetCommunityByID(p.CommunityID); err != nil {
		zap.L().Error("community not found", zap.Int64("community_id", p.CommunityID))
		return ErrorCommunityNotFound
	}

	postID := snowflake.GenerateID()
	p.ID = postID

	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}

	if err := redis.CreatePost(postID, p.CommunityID); err != nil {
		zap.L().Warn("redis.CreatePost failed",
			zap.Int64("post_id", postID),
			zap.Error(err))
	}
	return
}

func GetPostByID(postID int64) (data *post.DetailPost, err error) {
	postInfo, err := mysql.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	username, err := mysql.GetUserNameByID(postInfo.AuthorID)
	if err != nil {
		return nil, err
	}

	community, err := mysql.GetCommunityByID(postInfo.CommunityID)
	if err != nil {
		return nil, err
	}

	data = &post.DetailPost{
		Post:          postInfo,
		AuthorName:    username,
		CommunityName: community.CommunityName,
	}
	return data, nil
}

func GetPostList(page, size int) (data []*post.DetailPost, err error) {
	postList, err := mysql.GetPostList(page, size)
	data = make([]*post.DetailPost, 0, len(postList))
	for _, p := range postList {
		username, err := mysql.GetUserNameByID(p.AuthorID)
		if err != nil {
			zap.L().Warn("get username failed, skip post",
				zap.Int64("post_id", p.ID),
				zap.Int64("author_id", p.AuthorID),
				zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityByID(p.CommunityID)
		if err != nil {
			zap.L().Warn("get community failed, skip post",
				zap.Int64("post_id", p.ID),
				zap.Int64("community_id", p.CommunityID),
				zap.Error(err))
			continue
		}
		detailPost := &post.DetailPost{
			Post:          p,
			AuthorName:    username,
			CommunityName: community.CommunityName,
		}
		data = append(data, detailPost)
	}
	return
}

func GetPostListByOrder(p *param.PostsPageParams) (data []*post.DetailPost, err error) {
	var ids []string
	if p.CommunityID > 0 {
		_, err = mysql.GetCommunityByID(p.CommunityID)
		if err != nil {
			return nil, ErrorCommunityNotFound
		}
		ids, err = redis.GetPostsIdInorderWCommu(p)
	} else {
		ids, err = redis.GetPostsIdInorder(p)
	}
	if err != nil {
		zap.L().Error("redis.GetPostsIdInorder failed",
			zap.Error(err))
		return
	}
	if len(ids) == 0 {
		data = make([]*post.DetailPost, 0)
		return
	}
	// mysql查询获取帖子
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostListByIds failed",
			zap.Error(err))
		return
	}
	// 获取票数
	voteNums, err := redis.GetPostVotesByIds(ids)
	if err != nil {
		zap.L().Error("redis.GetPostVoteNums failed",
			zap.Error(err))
		return
	}
	// 查询作者、分区名字并加入
	data = make([]*post.DetailPost, 0, len(posts))
	for i, p := range posts {
		username, err := mysql.GetUserNameByID(p.AuthorID)
		if err != nil {
			zap.L().Warn("get username failed, skip post",
				zap.Int64("post_id", p.ID),
				zap.Int64("author_id", p.AuthorID),
				zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityByID(p.CommunityID)
		if err != nil {
			zap.L().Warn("get community failed, skip post",
				zap.Int64("post_id", p.ID),
				zap.Int64("community_id", p.CommunityID),
				zap.Error(err))
			continue
		}
		// // 获取票数并拼接
		// voteNum, err := redis.GetPostVoteNum(p.ID)
		// if err != nil {
		// 	zap.L().Warn("get vote num failed, skip post",
		// 		zap.Int64("post_id", p.ID),
		// 		zap.Error(err))
		// 	continue
		// }
		detailPost := &post.DetailPost{
			VoteNum:       voteNums[i],
			Post:          p,
			AuthorName:    username,
			CommunityName: community.CommunityName,
		}
		data = append(data, detailPost)
	}
	return
}

// func GetPostListByOrderAndCommu(p *param.PostsPageParams) (data []*post.DetailPost, err error) {
// 	// redis获取目标ids
// 	ids, err := redis.GetPostsIdInorderWCommu(p)
// 	if err != nil {
// 		zap.L().Error("redis.GetPostsIdInorder failed",
// 			zap.Error(err))
// 		return
// 	}
// 	if len(ids) == 0 {
// 		data = make([]*post.DetailPost, 0)
// 		return
// 	}
// 	// mysql查询获取帖子
// 	posts, err := mysql.GetPostListByIds(ids)
// 	if err != nil {
// 		zap.L().Error("mysql.GetPostListByIds failed",
// 			zap.Error(err))
// 		return
// 	}
// 	// 获取票数
// 	voteNums, err := redis.GetPostVotesByIds(ids)
// 	if err != nil {
// 		zap.L().Error("redis.GetPostVoteNums failed",
// 			zap.Error(err))
// 		return
// 	}
// 	// 查询作者、分区名字并加入
// 	data = make([]*post.DetailPost, 0, len(posts))
// 	for i, p := range posts {
// 		username, err := mysql.GetUserNameByID(p.AuthorID)
// 		if err != nil {
// 			zap.L().Warn("get username failed, skip post",
// 				zap.Int64("post_id", p.ID),
// 				zap.Int64("author_id", p.AuthorID),
// 				zap.Error(err))
// 			continue
// 		}
// 		community, err := mysql.GetCommunityByID(p.CommunityID)
// 		if err != nil {
// 			zap.L().Warn("get community failed, skip post",
// 				zap.Int64("post_id", p.ID),
// 				zap.Int64("community_id", p.CommunityID),
// 				zap.Error(err))
// 			continue
// 		}
// 		// // 获取票数并拼接
// 		// voteNum, err := redis.GetPostVoteNum(p.ID)
// 		// if err != nil {
// 		// 	zap.L().Warn("get vote num failed, skip post",
// 		// 		zap.Int64("post_id", p.ID),
// 		// 		zap.Error(err))
// 		// 	continue
// 		// }
// 		detailPost := &post.DetailPost{
// 			VoteNum:       voteNums[i],
// 			Post:          p,
// 			AuthorName:    username,
// 			CommunityName: community.CommunityName,
// 		}
// 		data = append(data, detailPost)
// 	}
// 	return
// }
