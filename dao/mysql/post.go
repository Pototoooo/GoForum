package mysql

import (
	"strings"

	"bluebell/models/post"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func CreatePost(p *post.Post) error {
	sqlStr := `insert into post (post_id, title, content, author_id, community_id)
	values(?,?,?,?,?)`
	// 执行sql
	_, err := db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	if err != nil {
		zap.L().Error("create post failed", zap.Error(err))
		return err
	}
	return nil
}

func GetPostByID(postID int64) (data *post.Post, err error) {
	data = new(post.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time, update_time 
	from post where post_id = ?`
	err = db.Get(data, sqlStr, postID)
	if err != nil {
		zap.L().Error("get post by id failed", zap.Error(err))
		return nil, err
	}
	return data, nil
}

func GetPostList(page, size int) (data []*post.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id,status, create_time, update_time 
	from post order by create_time desc limit ?,?`
	err = db.Select(&data, sqlStr, page-1, size)
	if err != nil {
		zap.L().Error("get post list failed", zap.Error(err))
		return nil, err
	}
	return data, nil
}

func GetPostListByIds(Ids []string) (postList []*post.Post, err error) {
	if len(Ids) == 0 {
		return
	}
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id,?)
	`
	query, args, err := sqlx.In(sqlStr, Ids, strings.Join(Ids, ","))
	if err != nil {
		return
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}
