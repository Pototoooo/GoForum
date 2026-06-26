package mysql

import (
	"database/sql"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type PostVote struct {
	PostID    int64 `db:"post_id"`
	UserID    int64 `db:"user_id"`
	Direction int   `db:"direction"`
}

func ensurePostVoteTable() error {
	sqlStr := `CREATE TABLE IF NOT EXISTS post_vote (
		id bigint(20) NOT NULL AUTO_INCREMENT,
		post_id bigint(20) NOT NULL COMMENT '帖子id',
		user_id bigint(20) NOT NULL COMMENT '投票用户id',
		direction tinyint(4) NOT NULL COMMENT '投票方向：1赞成，-1反对',
		create_time timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
		update_time timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
		PRIMARY KEY (id),
		UNIQUE KEY idx_user_post (user_id, post_id) USING BTREE,
		KEY idx_post_id (post_id) USING BTREE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`
	_, err := db.Exec(sqlStr)
	return err
}

func GetPostVoteDirection(userID, postID int64) (direction int, err error) {
	sqlStr := `select direction from post_vote where user_id = ? and post_id = ?`
	err = db.Get(&direction, sqlStr, userID, postID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return
}

func UpsertPostVote(userID, postID int64, direction int) error {
	sqlStr := `insert into post_vote (user_id, post_id, direction)
	values (?, ?, ?)
	on duplicate key update direction = values(direction)`
	_, err := db.Exec(sqlStr, userID, postID, direction)
	return err
}

func DeletePostVote(userID, postID int64) error {
	sqlStr := `delete from post_vote where user_id = ? and post_id = ?`
	_, err := db.Exec(sqlStr, userID, postID)
	return err
}

func GetPostVoteNumsByIds(ids []string) (map[int64]int, error) {
	result := make(map[int64]int, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	query, args, err := sqlx.In(`
		select post_id, count(*) as vote_num
		from post_vote
		where direction = 1 and post_id in (?)
		group by post_id`, ids)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)

	var rows []struct {
		PostID  int64 `db:"post_id"`
		VoteNum int   `db:"vote_num"`
	}
	if err := db.Select(&rows, query, args...); err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.PostID] = row.VoteNum
	}
	for _, id := range ids {
		postID, err := strconv.ParseInt(id, 10, 64)
		if err == nil {
			if _, ok := result[postID]; !ok {
				result[postID] = 0
			}
		}
	}
	return result, nil
}

func GetAllPostVotes() ([]PostVote, error) {
	var votes []PostVote
	sqlStr := `select post_id, user_id, direction from post_vote`
	err := db.Select(&votes, sqlStr)
	return votes, err
}
