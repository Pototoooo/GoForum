package route

import (
	"bluebell/controller"
	_ "bluebell/docs"
	"bluebell/logger"
	"bluebell/middlewire"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup() (r *gin.Engine) {
	r = gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 静态文件服务
	r.Static("/static", "./bluebell_frontend/dist/static")
	r.GET("/", func(c *gin.Context) {
		c.File("./bluebell_frontend/dist/index.html")
	})

	// API v1 路由组（前端 bluebell_frontend 使用）
	v1 := r.Group("/api/v1")
	{
		v1.POST("/signup", controller.RegisterHandler)
		v1.POST("/login", controller.LoginHandler)
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.GET("/posts2", controller.GetOrderedPosts)
		v1.GET("/post/:id", controller.GetDetailPostHandler)

		v1.Use(middlewire.RateLimiter(50, 100), middlewire.JWTAuthMidWire())
		{
			v1.POST("/post", controller.PostHandler)
			v1.POST("/vote", controller.VoteHandler)
		}
	}

	// SPA 历史模式路由：所有非 API 路径都返回 index.html
	r.NoRoute(func(c *gin.Context) {
		c.File("./bluebell_frontend/dist/index.html")
	})
	return
}

// 注册业务的路由
// func RegisterRoutes(r *gin.Engine) {
// 	// 注册路由
// 	r.POST("/register", register.RegisterHandler)
// }
