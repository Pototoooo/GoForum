package logic

import (
	"time"

	"GoForum/dao/mysql"
	"GoForum/dao/redis"

	"go.uber.org/zap"
)

func VoteHandler(userID int64, direction int, postID int64) (err error) {
	postInfo, err := mysql.GetPostByID(postID)
	if err != nil {
		return err
	}
	if time.Since(postInfo.CreateTime) > redis.VoteExpireSeconds*time.Second {
		return redis.ErrorVoteTimeExpire
	}

	oldDirection, err := mysql.GetPostVoteDirection(userID, postID)
	if err != nil {
		return err
	}
	if oldDirection == direction {
		if direction == 0 {
			return nil
		}
		return redis.ErrorVotedRepeated
	}

	if direction == 0 {
		err = mysql.DeletePostVote(userID, postID)
	} else {
		err = mysql.UpsertPostVote(userID, postID, direction)
	}
	if err != nil {
		return err
	}

	if err := redis.ApplyVoteChange(userID, postID, oldDirection, direction); err != nil {
		zap.L().Warn("redis.ApplyVoteChange failed",
			zap.Int64("user_id", userID),
			zap.Int64("post_id", postID),
			zap.Int("old_direction", oldDirection),
			zap.Int("new_direction", direction),
			zap.Error(err))
	}
	return nil
}
