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
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "登录",
	})
}

// Login 处理登录提交
func (h *AuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	cfg := config.Get()

	// 调试
	if cfg == nil {
		c.String(500, "config is nil!")
		return
	}

	if username == cfg.Admin.Username && password == cfg.Admin.Password {
		// 登录成功，设置Cookie
		c.SetCookie("admin_session", username, 86400, "/admin", "", false, true)
		// 返回重定向
		c.Header("Location", "/admin/")
		c.Status(302)
		return
	}

	// 登录失败 - 需要手动设置状态
	c.Status(200)
	c.HTML(200, "login.html", gin.H{
		"title": "登录",
		"error": "用户名或密码错误",
	})
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "admin_session",
		Value:  "",
		Path:   "/admin",
		MaxAge: -1,
	})
	c.Redirect(http.StatusFound, "/admin/login")
}

// AuthRequired 检查是否已登录
func (h *AuthHandler) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("admin_session")
		if err != nil || cookie == nil || cookie.Value == "" {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
