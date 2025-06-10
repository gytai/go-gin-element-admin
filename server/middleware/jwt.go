package middleware

import (
	"net/http"
	"server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从header中获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "请求未携带token，无权限访问",
			})
			c.Abort()
			return
		}

		// 移除Bearer前缀
		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}

		// 解析token
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "无效的token",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到context中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("authorityId", claims.AuthorityId)

		c.Next()
	}
}

// GetCurrentUser 从context中获取当前用户信息
func GetCurrentUser(c *gin.Context) (userID uint, username string, authorityId uint, exists bool) {
	userIDInterface, exists1 := c.Get("userID")
	usernameInterface, exists2 := c.Get("username")
	authorityIdInterface, exists3 := c.Get("authorityId")

	if !exists1 || !exists2 || !exists3 {
		return 0, "", 0, false
	}

	userID, ok1 := userIDInterface.(uint)
	username, ok2 := usernameInterface.(string)
	authorityId, ok3 := authorityIdInterface.(uint)

	if !ok1 || !ok2 || !ok3 {
		return 0, "", 0, false
	}

	return userID, username, authorityId, true
}
