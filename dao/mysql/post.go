package mysql

import (
	"strings"
	"time"

	"GoForum/models/param"
	"GoForum/models/post"
	"GoForum/pkg"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PostIndexRow struct {
	PostID      int64     `db:"post_id"`
	CommunityID int64     `db:"community_id"`
	CreateTime  time.Time `db:"create_time"`
	Score       int       `db:"score"`
}

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

func GetPostListByOrderWithVotes(p *param.PostsPageParams) (postList []*post.Post, voteNums map[int64]int, err error) {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size < 1 {
		p.Size = 10
	}
	offset := (p.Page - 1) * p.Size

	where := ""
	args := []interface{}{}
	if p.CommunityID > 0 {
		where = "where p.community_id = ?"
		args = append(args, p.CommunityID)
	}

	orderBy := "p.create_time desc"
	if p.Order == pkg.OrderScore {
		orderBy = "score desc, p.create_time desc"
	}

	sqlStr := `select p.post_id, p.title, p.content, p.author_id, p.community_id, p.status,
		p.create_time, p.update_time,
		coalesce(sum(case when pv.direction = 1 then 1 else 0 end), 0) as vote_num,
		coalesce(sum(pv.direction), 0) as score
	from post p
	left join post_vote pv on p.post_id = pv.post_id
	` + where + `
	group by p.id, p.post_id, p.title, p.content, p.author_id, p.community_id, p.status, p.create_time, p.update_time
	order by ` + orderBy + `
	limit ?, ?`
	args = append(args, offset, p.Size)

	var rows []struct {
		ID          int64     `db:"post_id"`
		Title       string    `db:"title"`
		Content     string    `db:"content"`
		AuthorID    int64     `db:"author_id"`
		CommunityID int64     `db:"community_id"`
		Status      int32     `db:"status"`
		CreateTime  time.Time `db:"create_time"`
		UpdateTime  time.Time `db:"update_time"`
		VoteNum     int       `db:"vote_num"`
		Score       int       `db:"score"`
	}
	err = db.Select(&rows, sqlStr, args...)
	if err != nil {
		return nil, nil, err
	}

	voteNums = make(map[int64]int, len(rows))
	postList = make([]*post.Post, 0, len(rows))
	for _, row := range rows {
		postList = append(postList, &post.Post{
			ID:          row.ID,
			Title:       row.Title,
			Content:     row.Content,
			AuthorID:    row.AuthorID,
			CommunityID: row.CommunityID,
			Status:      row.Status,
			CreateTime:  row.CreateTime,
			UpdateTime:  row.UpdateTime,
		})
		voteNums[row.ID] = row.VoteNum
	}
	return
}

func GetAllPostIndexRows() ([]PostIndexRow, error) {
	var rows []PostIndexRow
	sqlStr := `select p.post_id, p.community_id, p.create_time,
		coalesce(sum(pv.direction), 0) as score
	from post p
	left join post_vote pv on p.post_id = pv.post_id
	group by p.id, p.post_id, p.community_id, p.create_time
	order by p.id asc`
	err := db.Select(&rows, sqlStr)
	return rows, err
}
