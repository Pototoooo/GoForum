package middlewire

import (
	"GoForum/pkg/Code"
	"GoForum/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimiter(r float64, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(r), b)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			response.ResponseError(c, Code.CodeTooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}
}
