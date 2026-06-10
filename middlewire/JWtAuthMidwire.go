package middlewire

import (
	"strings"

	"bluebell/pkg/Code"
	"bluebell/pkg/jwt"
	"bluebell/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	UserID   = "UserId"
	UserName = "UserName"
)

func JWTAuthMidWire() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ResponseError(c, Code.CodeTokenError)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.ResponseError(c, Code.CodeTokenError)
			c.Abort()
			return
		}
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			zap.L().Error("JWT ParseToken failed", zap.Error(err))
			response.ResponseError(c, Code.CodeTokenError)
			c.Abort()
			return
		}
		c.Set(UserID, mc.UserID)
		c.Set(UserName, mc.Username)
		c.Next()
	}
}
