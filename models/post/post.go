package post

import (
	"time"
)

type Post struct {
	// 内存对齐
	Title       string    `json:"title" db:"title" required:"true"`
	Content     string    `json:"content" db:"content" required:"true"`
	ID          int64     `json:"id,string" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" required:"true"`
	Status      int32     `json:"status" db:"status"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
	UpdateTime  time.Time `json:"update_time" db:"update_time"`
}
type DetailPost struct {
	AuthorName    string `json:"author_name"`
	CommunityName string `json:"community_name"`
	VoteNum       int `json:"vote_num"`
	*Post
}
