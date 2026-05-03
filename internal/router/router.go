package router

import (
	"html/template"
	"log"
	"os"
	"strconv"
	"time"

	"Goblog/internal/config"
	"Goblog/internal/handler/admin"
	"Goblog/internal/handler/front"
	"Goblog/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupTemplateFunctions 设置模板函数
func SetupTemplateFunctions(engine *gin.Engine) {
	funcMap := template.FuncMap{
		"add":       func(a, b int) int { return a + b },
		"sub":       func(a, b int) int { return a - b },
		"eq":        func(a, b interface{}) bool { return a == b },
		"gt":        func(a, b int) bool { return a > b },
		"lt":        func(a, b int) bool { return a < b },
		"not":       func(b bool) bool { return !b },
		"safe":      func(str string) template.HTML { return template.HTML(str) },
		"strToUint": func(s string) uint { n, _ := strconv.ParseUint(s, 10, 64); return uint(n) },
		"seq": func(n int) []int {
			s := make([]int, n)
			for i := 0; i < n; i++ {
				s[i] = i + 1
			}
			return s
		},
		"pageRange": func(current, total int) []int {
			var pages []int
			// 最多显示5个页码
			if total <= 5 {
				for i := 1; i <= total; i++ {
					pages = append(pages, i)
				}
			} else {
				// 始终显示第一页
				pages = append(pages, 1)
				// 计算中间页码范围
				start := current - 1
				end := current + 1
				if start <= 2 {
					start = 2
					end = 4
				}
				if end >= total-1 {
					start = total - 3
					end = total - 1
				}
				if start > 2 {
					pages = append(pages, -1) // 省略号
				}
				for i := start; i <= end; i++ {
					pages = append(pages, i)
				}
				if end < total-1 {
					pages = append(pages, -1) // 省略号
				}
				// 始终显示最后一页
				pages = append(pages, total)
			}
			return pages
		},
		"formatTime": func(unix interface{}) string {
			if unix == nil {
				return ""
			}
			var ts int64
			switch v := unix.(type) {
			case int64:
				ts = v
			case int:
				ts = int64(v)
			case uint:
				ts = int64(v)
			case uint64:
				ts = int64(v)
			case float64:
				ts = int64(v)
			default:
				// 尝试直接转换
				return ""
			}
			if ts == 0 {
				return ""
			}
			t := time.Unix(ts, 0)
			return t.Format("2006-01-02")
		},
		"formatTs": func(ts int64) string {
			if ts == 0 {
				return ""
			}
			t := time.Unix(ts, 0)
			return t.Format("2006-01-02 15:04:05")
		},
		"markdown": func(content string) template.HTML {
			if content == "" {
				return ""
			}
			// 直接返回，服务端已渲染的 Markdown 不需要再次渲染
			return template.HTML(content)
		},
	}
	engine.SetFuncMap(funcMap)
}

// Setup 设置路由
func Setup(
	engine *gin.Engine,
	postService *service.PostService,
	columnService *service.ColumnService,
	commentService *service.CommentService,
	visitorService *service.VisitorService,
	postLikeService *service.PostLikeService,
	fileService service.FileService,
	devlogService *service.DevlogService,
) {
	// 获取配置
	cfg := config.Get()
	adminPath := cfg.Admin.Path

	// 初始化后台Handler
	authHandler := admin.NewAuthHandler()
	postHandler := admin.NewPostHandler(postService, columnService)
	columnHandler := admin.NewColumnHandler(columnService)
	commentHandler := admin.NewCommentHandler(commentService)
	uploadHandler := admin.NewUploadHandler(fileService)
	adminDevlogHandler := admin.NewDevlogHandler(devlogService)

	// 初始化前台Handler
	homeHandler := front.NewHomeHandler(postService, columnService, devlogService)
	frontColumnHandler := front.NewColumnHandler(postService, columnService, commentService)
	messageHandler := front.NewMessageHandler(postService, commentService, visitorService)
	frontPostHandler := front.NewPostHandler(postService, commentService, visitorService, postLikeService)
	devlogHandler := front.NewDevlogHandler(devlogService)

	// 设置模板函数（必须在加载模板之前）
	SetupTemplateFunctions(engine)

	// 设置HTML模板（后台+前台）
	engine.LoadHTMLGlob("./web/templates/**/*.html")

	// 设置后台路径到 Gin context，供模板使用
	engine.Use(func(c *gin.Context) {
		c.Set("adminPath", adminPath)
		c.Next()
	})

	// ============ 前台路由 ============
	// 搜索API - 放在前面避免被其他路由匹配
	engine.GET("/api/search", homeHandler.Search)

	// 文件上传 API - 不需要认证（公共路由）
	engine.POST("/api/upload", uploadHandler.Upload)

	engine.GET("/", homeHandler.Index)
	engine.GET("/devlog", devlogHandler.Index)
	engine.GET("/column", frontColumnHandler.Index)
	engine.GET("/column/:slug", frontColumnHandler.List)
	engine.GET("/post/:slug", frontColumnHandler.Post)
	engine.GET("/message", messageHandler.Index)
	engine.POST("/message", messageHandler.Submit)

	// ============ 访客API ============
	engine.GET("/api/visitor/check", messageHandler.Check)
	engine.GET("/api/visitor/check-exist", messageHandler.CheckExist)
	engine.POST("/api/visitor/register", messageHandler.Register)
	engine.POST("/api/visitor/update", messageHandler.UpdateInfo)
	engine.DELETE("/api/visitor/all", messageHandler.DeleteAllVisitors)

	// ============ 评论API ============
	engine.POST("/api/comment", messageHandler.SubmitPostComment)

	// ============ 文章API ============
	engine.POST("/api/post/:id/like", frontPostHandler.Like)
	engine.GET("/api/post/:id/check-like", frontPostHandler.CheckLike)

	// ============ 后台公开路由（登录/登出/上传）===========
	adminPublic := engine.Group(adminPath)
	{
		adminPublic.GET("/login", authHandler.LoginPage)
		adminPublic.POST("/login", authHandler.Login)
		adminPublic.GET("/logout", authHandler.Logout)
		
		// 文件上传 - 不需要认证（能访问后台就能上传）
		adminPublic.POST("/upload", uploadHandler.Upload)
		adminPublic.POST("/upload/cover", uploadHandler.UploadCover)
	}

	// ============ 第二组：需要认证的路由 ============
	adminPrivate := engine.Group(adminPath)
	{
		// 使用认证中间件
		adminPrivate.Use(authHandler.AuthRequired())

// 仪表盘（带统计数据）
	adminPrivate.GET("/", func(c *gin.Context) {
		postsTotal, likesTotal, _, _ := postService.GetStats()
		columns, _ := columnService.GetAll()
		_, messageTotal, _ := commentService.GetAll(1, 1)
		c.HTML(200, "dashboard.html", gin.H{
			"title":        "后台管理",
			"PostsTotal":   postsTotal,
			"LikesTotal":   likesTotal,
			"ColumnsTotal": len(columns),
			"CommentTotal": messageTotal,
			"adminPath":   adminPath,
		})
	})

		// 文章管理
		adminPrivate.GET("/posts", postHandler.List)
		adminPrivate.GET("/posts/new", postHandler.Edit)
		adminPrivate.GET("/posts/:id", postHandler.Edit)
		adminPrivate.POST("/posts/save", postHandler.Save)
		adminPrivate.DELETE("/posts/:id", postHandler.Delete)
		adminPrivate.POST("/posts/:id/publish", postHandler.Publish)
		adminPrivate.POST("/posts/:id/migrate", postHandler.Migrate)
		adminPrivate.POST("/posts/batch-migrate", postHandler.BatchMigrate)

		// 草稿箱
		adminPrivate.GET("/drafts", postHandler.Drafts)

		// 专栏管理
		// 文章评论管理（PostID > 0）- 必须放在 columns/:id 之前
		adminPrivate.GET("/article-comments", commentHandler.ArticleComments)
		adminPrivate.POST("/article-comments/:id/approve", commentHandler.Approve)
		adminPrivate.POST("/article-comments/:id/reject", commentHandler.Reject)
		adminPrivate.POST("/article-comments/batch-approve", commentHandler.BatchApprove)
		adminPrivate.POST("/article-comments/batch-reject", commentHandler.BatchReject)
		adminPrivate.DELETE("/article-comments/:id", commentHandler.Delete)

		// 留言管理（PostID = 0）
		adminPrivate.GET("/comments", commentHandler.List)
		adminPrivate.POST("/comments/:id/approve", commentHandler.Approve)
		adminPrivate.POST("/comments/:id/reject", commentHandler.Reject)
		adminPrivate.POST("/comments/batch-approve", commentHandler.BatchApprove)
		adminPrivate.POST("/comments/batch-reject", commentHandler.BatchReject)
		adminPrivate.DELETE("/comments/:id", commentHandler.Delete)

		adminPrivate.GET("/columns", columnHandler.List)
		adminPrivate.GET("/columns/new", columnHandler.Edit)
		adminPrivate.GET("/columns/:id", columnHandler.Edit)
		adminPrivate.POST("/columns/save", columnHandler.Save)
		adminPrivate.DELETE("/columns/:id", columnHandler.Delete)

		// 开发日志
		adminPrivate.GET("/devlogs", adminDevlogHandler.List)
		adminPrivate.GET("/devlogs/new", adminDevlogHandler.Edit)
		adminPrivate.GET("/devlogs/:id", adminDevlogHandler.Edit)
		adminPrivate.POST("/devlogs/save", adminDevlogHandler.Save)
		adminPrivate.DELETE("/devlogs/:id", adminDevlogHandler.Delete)
		adminPrivate.POST("/devlogs/:id/publish", adminDevlogHandler.Publish)
		adminPrivate.POST("/devlogs/:id/unpublish", adminDevlogHandler.Unpublish)

		// 文件上传
	}

	// 设置静态文件 - 使用项目根目录
	wd, _ := os.Getwd()
	staticPath := wd + "/web/static"

	log.Printf("[INFO] Static /static mapped to: %s", staticPath)
	engine.Static("/static", staticPath)
}
