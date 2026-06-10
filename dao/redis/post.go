package redis

import (
	"strconv"
	"time"

	"bluebell/models/param"
	"bluebell/pkg"

	"github.com/redis/go-redis/v9"
)

func CreatePost(postID int64, communityID int64) error {
	now := float64(time.Now().Unix())
	pid := strconv.FormatInt(postID, 10)

	pipeline := rdb.TxPipeline()
	pipeline.ZAdd(ctx, pkg.KeyPostTime, redis.Z{
		Score:  now,
		Member: pid,
	})
	pipeline.ZAdd(ctx, pkg.KeyPostScore, redis.Z{
		Score:  0,
		Member: pid,
	})
	pipeline.SAdd(ctx, pkg.KeyCommunityPostPrefix+strconv.FormatInt(communityID, 10), pid)
	_, err := pipeline.Exec(ctx)
	return err
}

func GetPostsIdInorder(p *param.PostsPageParams) (Ids []string, err error) {
	key := pkg.KeyPostTime
	if p.Order == pkg.OrderScore {
		key = pkg.KeyPostScore
	}
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	return rdb.ZRevRange(ctx, key, int64(start), int64(end)).Result()
}

// GetPostVotesByIds 根据帖子ID列表查询每个帖子的赞成票数
func GetPostVotesByIds(ids []string) (data []int, err error) {
	p := rdb.Pipeline()
	data = make([]int, 0)
	for _, id := range ids {
		key := pkg.KeyPostVotedPrefix + id
		p.ZCount(ctx, key, "1", "1")
	}
	cmders, err := p.Exec(ctx)
	if err != nil {
		return nil, err
	}
	for _, cmder := range cmders {
		value := cmder.(*redis.IntCmd).Val()
		data = append(data, int(value))
	}
	return
}

func GetPostsIdInorderWCommu(p *param.PostsPageParams) (Ids []string, err error) {
	orderKey := pkg.KeyPostTime
	communityIdStr := strconv.FormatInt(p.CommunityID, 10)
	communityKey := pkg.KeyCommunityPostPrefix + communityIdStr
	if p.Order == pkg.OrderScore {
		orderKey = pkg.KeyPostScore
	}
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	tempKey := orderKey + ":" + communityIdStr

	if rdb.Exists(ctx, tempKey).Val() < 1 {
		err = rdb.ZInterStore(ctx, tempKey, &redis.ZStore{
			Keys:      []string{orderKey, communityKey},
			Aggregate: "MAX",
		}).Err()
		if err != nil {
			return nil, err
		}
		rdb.Expire(ctx, tempKey, 60*time.Second)
	}

	return rdb.ZRevRange(ctx, tempKey, int64(start), int64(end)).Result()
}
