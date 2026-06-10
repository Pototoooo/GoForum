package logic

import "bluebell/dao/redis"

func VoteHandler(userID int64, direction int, postID int64) (err error) {
	err = redis.VoteHandler(userID, direction, postID)
	return
}
