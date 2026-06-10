package redis

import (
	"context"
	"fmt"

	"bluebell/settings"

	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx context.Context
)

func Init() (err error) {
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", settings.Config.Redis.Host, settings.Config.Redis.Port),
		Password: settings.Config.Redis.Password,
		DB:       settings.Config.Redis.DB,
		PoolSize: settings.Config.Redis.PoolSize,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("can't connect to redis")
		return
	}
	return
}

func Close() {
	rdb.Close()
}
