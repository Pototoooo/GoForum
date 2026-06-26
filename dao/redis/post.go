package redis

import (
	"strconv"
	"time"

	"GoForum/models/param"
	"GoForum/pkg"

	"github.com/redis/go-redis/v9"
)

type PostIndexItem struct {
	PostID      int64
	CommunityID int64
	CreateTime  time.Time
	Score       int
}

type VoteItem struct {
	PostID    int64
	UserID    int64
	Direction int
}

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

func RebuildPostIndex(posts []PostIndexItem, votes []VoteItem) error {
	if err := clearPostIndex(); err != nil {
		return err
	}

	pipeline := rdb.TxPipeline()
	for _, p := range posts {
		postID := strconv.FormatInt(p.PostID, 10)
		communityID := strconv.FormatInt(p.CommunityID, 10)
		pipeline.ZAdd(ctx, pkg.KeyPostTime, redis.Z{
			Score:  float64(p.CreateTime.Unix()),
			Member: postID,
		})
		pipeline.ZAdd(ctx, pkg.KeyPostScore, redis.Z{
			Score:  float64(p.Score) * VoteScore,
			Member: postID,
		})
		pipeline.SAdd(ctx, pkg.KeyCommunityPostPrefix+communityID, postID)
	}
	for _, vote := range votes {
		if vote.Direction == 0 {
			continue
		}
		postID := strconv.FormatInt(vote.PostID, 10)
		userID := strconv.FormatInt(vote.UserID, 10)
		pipeline.ZAdd(ctx, pkg.KeyPostVotedPrefix+postID, redis.Z{
			Score:  float64(vote.Direction),
			Member: userID,
		})
	}
	_, err := pipeline.Exec(ctx)
	return err
}

func clearPostIndex() error {
	keys := []string{pkg.KeyPostTime, pkg.KeyPostScore}
	iter := rdb.Scan(ctx, 0, pkg.KeyPostVotedPrefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	iter = rdb.Scan(ctx, 0, pkg.KeyCommunityPostPrefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return rdb.Del(ctx, keys...).Err()
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
