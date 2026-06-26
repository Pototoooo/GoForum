package redis

import (
	"errors"
	"strconv"
	"time"

	"GoForum/pkg"

	"github.com/redis/go-redis/v9"
)

const (
	VoteExpireSeconds = 7 * 24 * 3600
	VoteScore         = 432.0
)

var (
	ErrorVoteTimeExpire = errors.New("已过投票时间")
	ErrorVotedRepeated  = errors.New("重复投票")
)

func VoteHandler(userID int64, direction int, postID int64) error {
	// 转string
	postIDStr := strconv.FormatInt(postID, 10)
	userIDStr := strconv.FormatInt(userID, 10)
	// 检查是否过期
	postTime := rdb.ZScore(ctx, pkg.KeyPostTime, postIDStr).Val()
	if float64(time.Now().Unix())-postTime > VoteExpireSeconds {
		return ErrorVoteTimeExpire
	}

	votedKey := pkg.KeyPostVotedPrefix + postIDStr
	// 获取之前的投票方向
	oldDirection := rdb.ZScore(ctx, votedKey, userIDStr).Val()

	if direction == 0 {
		if oldDirection == 0 {
			return nil
		}
		pipeline := rdb.TxPipeline()
		pipeline.ZRem(ctx, votedKey, userIDStr)
		pipeline.ZIncrBy(ctx, pkg.KeyPostScore, -oldDirection*VoteScore, postIDStr)
		_, err := pipeline.Exec(ctx)
		return err
	}

	if oldDirection == 0 {
		pipeline := rdb.TxPipeline()
		pipeline.ZAdd(ctx, votedKey, redis.Z{Score: float64(direction), Member: userIDStr})
		pipeline.ZIncrBy(ctx, pkg.KeyPostScore, float64(direction)*VoteScore, postIDStr)
		_, err := pipeline.Exec(ctx)
		return err
	}

	if int(oldDirection) == direction {
		return ErrorVotedRepeated
	}
	// 反转投票情况
	diff := float64(direction - int(oldDirection))
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(ctx, votedKey, diff, userIDStr)
	pipeline.ZIncrBy(ctx, pkg.KeyPostScore, diff*VoteScore, postIDStr)
	_, err := pipeline.Exec(ctx)
	return err
}

func ApplyVoteChange(userID int64, postID int64, oldDirection int, newDirection int) error {
	if oldDirection == newDirection {
		return nil
	}

	postIDStr := strconv.FormatInt(postID, 10)
	userIDStr := strconv.FormatInt(userID, 10)
	votedKey := pkg.KeyPostVotedPrefix + postIDStr
	diff := float64(newDirection - oldDirection)

	pipeline := rdb.TxPipeline()
	if newDirection == 0 {
		pipeline.ZRem(ctx, votedKey, userIDStr)
	} else {
		pipeline.ZAdd(ctx, votedKey, redis.Z{
			Score:  float64(newDirection),
			Member: userIDStr,
		})
	}
	pipeline.ZIncrBy(ctx, pkg.KeyPostScore, diff*VoteScore, postIDStr)
	_, err := pipeline.Exec(ctx)
	return err
}
