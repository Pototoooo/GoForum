package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"GoForum/dao/mysql"
	"GoForum/dao/redis"
	"GoForum/logger"
	"GoForum/pkg/snowflake"
	"GoForum/route"
	"GoForum/settings"

	"github.com/gin-contrib/pprof"
	"go.uber.org/zap"
)

// @title           GoForum API
// @version         1.0
// @description     GoForum 论坛项目 API 接口文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8084
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 JWT 认证，格式：Bearer {token}

func main() {
	// 初始化配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}

	// 初始化日志
	if err := logger.Init(settings.Config.Log.Console); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		os.Exit(1)
	}
	defer zap.L().Sync()

	// 初始化 MySQL
	if err := mysql.Init(); err != nil {
		zap.L().Fatal("init mysql failed", zap.Error(err))
	}
	defer mysql.Close()

	// 初始化 Redis
	if err := redis.Init(); err != nil {
		zap.L().Fatal("init redis failed", zap.Error(err))
	}
	defer redis.Close()

	// 初始化雪花算法
	if err := snowflake.Init(settings.Config.StartTime, settings.Config.MachineID); err != nil {
		zap.L().Fatal("init snowflake failed", zap.Error(err))
	}

	// 初始化路由
	r := route.Setup()
	pprof.Register(r)
	// 启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Config.Port),
		Handler: r,
	}

	go func() {
		zap.L().Info("server starting", zap.Int("port", settings.Config.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("server failed to start", zap.Error(err))
		}
	}()

	// 等待信号量，收到信号后关闭服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("服务关闭中")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("强制关机", zap.Error(err))
	}
	zap.L().Info("服务器已安全退出")
}
