package middleware

import (
	"Goblog/internal/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware 后台认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("admin_user_id")
		if !exists || userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// BasicAuthMiddleware HTTP Basic认证中间件
func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminPath := config.Get().Admin.Path

		// 登录页面不需要认证
		if c.Request.URL.Path == adminPath+"/login" || c.Request.URL.Path == adminPath+"/login/" {
			c.Next()
			return
		}

		// 非后台路径不需要认证
		if !strings.HasPrefix(c.Request.URL.Path, adminPath) {
			c.Next()
			return
		}

		username, password, hasAuth := c.Request.BasicAuth()
		if !hasAuth {
			c.Header("WWW-Authenticate", `Basic realm="Admin"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
			c.Abort()
			return
		}

		cfg := config.Get()
		if username != cfg.Admin.Username || password != cfg.Admin.Password {
			c.Header("WWW-Authenticate", `Basic realm="Admin"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证失败"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return gin.Logger()
}
