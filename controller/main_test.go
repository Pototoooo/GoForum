package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"GoForum/dao/mysql"
	"GoForum/dao/redis"
	"GoForum/pkg/snowflake"
	"GoForum/settings"
)

func TestMain(m *testing.M) {
	// 切换到项目根目录，使 config.yaml 路径解析正确
	projectRoot, err := findProjectRoot()
	if err != nil {
		fmt.Printf("找不到项目根目录: %v\n", err)
		fmt.Println("跳过依赖外部服务的测试")
		os.Exit(m.Run())
	}
	if err := os.Chdir(projectRoot); err != nil {
		fmt.Printf("切换目录失败: %v\n", err)
		os.Exit(m.Run())
	}

	// 初始化配置
	if err := settings.Init(); err != nil {
		fmt.Printf("初始化配置失败: %v\n", err)
		fmt.Println("跳过依赖外部服务的测试")
		os.Exit(m.Run())
	}

	// 初始化 MySQL
	if err := mysql.Init(); err != nil {
		fmt.Printf("初始化 MySQL 失败: %v\n", err)
		fmt.Println("跳过依赖 MySQL 的测试")
		os.Exit(m.Run())
	}
	defer mysql.Close()

	// 初始化 Redis
	if err := redis.Init(); err != nil {
		fmt.Printf("初始化 Redis 失败: %v\n", err)
		fmt.Println("跳过依赖 Redis 的测试")
		os.Exit(m.Run())
	}
	defer redis.Close()

	// 初始化雪花算法
	if err := snowflake.Init(settings.Config.StartTime, settings.Config.MachineID); err != nil {
		fmt.Printf("初始化 snowflake 失败: %v\n", err)
		os.Exit(m.Run())
	}

	// 运行所有测试
	code := m.Run()
	os.Exit(code)
}

// findProjectRoot 从当前目录向上查找包含 config.yaml 的目录
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "config.yaml")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("找不到 config.yaml")
		}
		dir = parent
	}
}
