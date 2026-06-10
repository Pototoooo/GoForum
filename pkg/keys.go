package pkg

const (
	// KeyPostInfoPrefix 帖子信息 Hash 的 key 前缀（完整 key: prefix + postID）
	KeyPostInfoPrefix = "bluebell:post:"
	// KeyPostTime 帖子创建时间有序集合（ZSet）
	KeyPostTime = "bluebell:post:time"
	// KeyPostScore 帖子分数有序集合（ZSet）
	KeyPostScore = "bluebell:post:score"
	// KeyPostVotedPrefix 每个帖子的投票用户集合前缀（完整 key: prefix + postID）
	KeyPostVotedPrefix = "bluebell:post:voted:"

	// KeyCommunityPostPrefix 社区帖子集合前缀（完整 key: prefix + communityID）
	KeyCommunityPostPrefix = "bluebell:community:"
)

// GetRedisKey 根据前缀和ID获取完整的Redis键
func GetRedisKey(prefix, id string) string {
	return prefix + id
}

// GetPostInfoKey 获取帖子信息 key
func GetPostInfoKey(postID string) string {
	return GetRedisKey(KeyPostInfoPrefix, postID)
}

// GetPostVotedKey 获取帖子投票 key
func GetPostVotedKey(postID string) string {
	return GetRedisKey(KeyPostVotedPrefix, postID)
}

// GetCommunityPostKey 获取社区帖子集合 key
func GetCommunityPostKey(communityID string) string {
	return GetRedisKey(KeyCommunityPostPrefix, communityID)
}
