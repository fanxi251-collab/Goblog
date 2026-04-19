package router

import (
	"Goblog/internal/config"
	"Goblog/internal/handler/admin"
	"Goblog/internal/handler/front"
	"Goblog/internal/service"

	"github.com/gin-gonic/gin"
	"html/template"
)

// SetupTemplateFunctions 设置模板函数
func SetupTemplateFunctions(engine *gin.Engine) {
	funcMap := template.FuncMap{
		"add":  func(a, b int) int { return a + b },
		"sub":  func(a, b int) int { return a - b },
		"eq":   func(a, b interface{}) bool { return a == b },
		"gt":   func(a, b int) bool { return a > b },
		"lt":   func(a, b int) bool { return a < b },
		"not":  func(b bool) bool { return !b },
		"safe": func(str string) template.HTML { return template.HTML(str) },
	}
	engine.SetFuncMap(funcMap)
}

// Setup 设置路由
func Setup(
	engine *gin.Engine,
	postService *service.PostService,
	columnService *service.ColumnService,
	commentService *service.CommentService,
	fileService service.FileService,
) {
	// 获取配置
	cfg := config.Get()
	adminPath := cfg.Admin.Path

	// 初始化后台Handler
	authHandler := admin.NewAuthHandler()
	postHandler := admin.NewPostHandler(postService, columnService)
	columnHandler := admin.NewColumnHandler(columnService)
	uploadHandler := admin.NewUploadHandler(fileService)

	// 初始化前台Handler
	homeHandler := front.NewHomeHandler(postService, columnService)
	frontColumnHandler := front.NewColumnHandler(postService, columnService)
	messageHandler := front.NewMessageHandler(commentService)

	// 设置模板函数（必须在加载模板之前）
	SetupTemplateFunctions(engine)

	// 设置HTML模板（后台+前台）
	engine.LoadHTMLGlob("./web/templates/**/*.html")

	// ============ 前台路由 ============
	engine.GET("/", homeHandler.Index)
	engine.GET("/about", homeHandler.About)
	engine.GET("/column", frontColumnHandler.Index)
	engine.GET("/column/:slug", frontColumnHandler.List)
	engine.GET("/post/:slug", frontColumnHandler.Post)
	engine.GET("/message", messageHandler.Index)
	engine.POST("/message", messageHandler.Submit)

	// ============ 后台公开路由（登录/登出）===========
	adminPublic := engine.Group(adminPath)
	{
		adminPublic.GET("/login", authHandler.LoginPage)
		adminPublic.POST("/login", authHandler.Login)
		adminPublic.GET("/logout", authHandler.Logout)
	}

	// ============ 第二组：需要认证的路由 ============
	adminPrivate := engine.Group(adminPath)
	{
		// 使用认证中间件
		adminPrivate.Use(authHandler.AuthRequired())

		// 仪表盘（带统计数据）
		adminPrivate.GET("/", func(c *gin.Context) {
			_, postsTotal, _ := postService.GetAll(1, 1)
			columns, _ := columnService.GetAll()
			_, commentsTotal, _ := commentService.GetAll(1, 1)
			println("DEBUG postsTotal:", postsTotal)
			println("DEBUG columns:", columns)
			println("DEBUG commentsTotal:", commentsTotal)
			c.HTML(200, "dashboard.html", gin.H{
				"title":        "后台管理",
				"PostsTotal":   postsTotal,
				"ColumnsTotal": len(columns),
				"CommentTotal": commentsTotal,
			})
		})

		// 文章管理
		adminPrivate.GET("/posts", postHandler.List)
		adminPrivate.GET("/posts/new", postHandler.Edit)
		adminPrivate.GET("/posts/:id", postHandler.Edit)
		adminPrivate.POST("/posts/save", postHandler.Save)
		adminPrivate.DELETE("/posts/:id", postHandler.Delete)
		adminPrivate.POST("/posts/:id/publish", postHandler.Publish)

		// 草稿箱
		adminPrivate.GET("/drafts", postHandler.Drafts)

		// 专栏管理
		adminPrivate.GET("/columns", columnHandler.List)
		adminPrivate.GET("/columns/new", columnHandler.Edit)
		adminPrivate.GET("/columns/:id", columnHandler.Edit)
		adminPrivate.POST("/columns/save", columnHandler.Save)
		adminPrivate.DELETE("/columns/:id", columnHandler.Delete)

		// 文件上传
		adminPrivate.POST("/upload", uploadHandler.Upload)
	}

	// 设置静态文件
	engine.Static("/static", "./web/static")
}
