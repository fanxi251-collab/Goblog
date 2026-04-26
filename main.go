package main

import (
	"Goblog/internal/config"
	"Goblog/internal/model"
	"Goblog/internal/repository"
	"Goblog/internal/router"
	"Goblog/internal/service"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 调试：打印配置
	log.Printf("后台路径: %s", cfg.Admin.Path)
	log.Printf("后台用户名: %s", cfg.Admin.Username)

	// 设置XSS开关
	model.SetXSSEnabled(cfg.XSS.Enabled)

	// 2. 创建数据库目录
	dbPath := cfg.Database.Path
	if dir := filepath.Dir(dbPath); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("创建数据库目录失败: %v", err)
		}
	}

	// 3. 连接数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 4. 自动迁移
	if err := db.AutoMigrate(
		&model.User{},
		&model.Column{},
		&model.Post{},
		&model.Comment{},
		&model.Visitor{},
		&model.PostLike{},
		&model.Devlog{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 5. 创建默认管理员（如果不存在）
	userRepo := repository.NewUserRepository(db)
	if _, err := userRepo.GetByUsername(cfg.Admin.Username); err != nil {
		defaultUser := &model.User{
			Username: cfg.Admin.Username,
			Password: cfg.Admin.Password,
			Nickname: "管理员",
		}
		if err := userRepo.Create(defaultUser); err != nil {
			log.Printf("创建默认用户失败: %v", err)
		}
	}

	// 6. 初始化Repository层
	postRepo := repository.NewPostRepository(db)
	columnRepo := repository.NewColumnRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	visitorRepo := repository.NewVisitorRepository(db)
	postLikeRepo := repository.NewPostLikeRepository(db)
	devlogRepo := repository.NewDevlogRepository(db)

	// 7. 初始化Service层
	postService := service.NewPostService(postRepo)
	columnService := service.NewColumnService(columnRepo)
	commentService := service.NewCommentService(commentRepo)
	visitorService := service.NewVisitorService(visitorRepo, commentRepo)
	postLikeService := service.NewPostLikeService(postLikeRepo, postRepo)
	devlogService := service.NewDevlogService(devlogRepo)
	fileService := service.NewLocalFileService()

	// 8. 创建Gin引擎
	engine := gin.Default()

	// 配置路由
	router.Setup(engine, postService, columnService, commentService, visitorService, postLikeService, fileService, devlogService)

	// 9. 启动服务
	addr := cfg.Address()
	log.Printf("服务器启动成功: http://%s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}

	// 防止编译提前退出
	log.Println("按 Ctrl+C 停止服务")
}
