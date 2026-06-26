package main

import (
	"fmt"
	"os"

	"GoForum/dao/mysql"
	"GoForum/dao/redis"
	"GoForum/logic"
	"GoForum/settings"
)

func main() {
	if err := settings.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "init settings failed: %v\n", err)
		os.Exit(1)
	}
	if err := mysql.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "init mysql failed: %v\n", err)
		os.Exit(1)
	}
	defer mysql.Close()

	if err := redis.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "init redis failed: %v\n", err)
		os.Exit(1)
	}
	defer redis.Close()

	if err := logic.RebuildRedisPostIndex(); err != nil {
		fmt.Fprintf(os.Stderr, "rebuild redis post index failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("redis post index rebuilt successfully")
}
