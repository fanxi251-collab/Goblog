package admin

import (
	"Goblog/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// LoginPage 登录页面（GET）
func (h *AuthHandler) LoginPage(c *gin.Context) {
	cfg := config.Get()
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":      "登录",
		"adminPath": cfg.Admin.Path,
	})
}

// Login 处理登录提交
func (h *AuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	cfg := config.Get()

	if cfg == nil {
		c.String(500, "config is nil!")
		return
	}

	if username == cfg.Admin.Username && password == cfg.Admin.Password {
		// 登录成功，设置Cookie
		c.SetCookie("admin_session", username, 86400, cfg.Admin.Path, "", false, true)
		// 返回重定向到后台首页
		c.Redirect(http.StatusFound, cfg.Admin.Path+"/")
		return
	}

	// 登录失败
	c.HTML(http.StatusUnauthorized, "login.html", gin.H{
		"title":     "登录",
		"adminPath": cfg.Admin.Path,
		"error":     "用户名或密码错误",
	})
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	cfg := config.Get()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "admin_session",
		Value:  "",
		Path:   cfg.Admin.Path,
		MaxAge: -1,
	})
	c.Redirect(http.StatusFound, cfg.Admin.Path+"/login")
}

// AuthRequired 检查是否已登录
func (h *AuthHandler) AuthRequired() gin.HandlerFunc {
	cfg := config.Get()
	loginURL := cfg.Admin.Path + "/login"
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("admin_session")
		if err != nil || cookie == nil || cookie.Value == "" {
			// 如果是 AJAX 请求，返回 JSON 而不是重定向
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
				c.Abort()
				return
			}
			c.Redirect(http.StatusFound, loginURL)
			c.Abort()
			return
		}
		c.Next()
	}
}
