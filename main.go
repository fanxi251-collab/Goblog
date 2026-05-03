package main

import (
	"Goblog/internal/config"
	"Goblog/internal/model"
	"Goblog/internal/repository"
	"Goblog/internal/router"
	"Goblog/internal/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
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

	// 2. 连接 PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	log.Printf("使用 PostgreSQL 数据库: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// 3. 自动迁移
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

	// 4. 创建默认管理员（如果不存在）
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

	// 5. 初始化Repository层
	postRepo := repository.NewPostRepository(db)
	columnRepo := repository.NewColumnRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	visitorRepo := repository.NewVisitorRepository(db)
	postLikeRepo := repository.NewPostLikeRepository(db)
	devlogRepo := repository.NewDevlogRepository(db)

	// 6. 初始化Service层
	postService := service.NewPostService(postRepo)
	columnService := service.NewColumnService(columnRepo)
	commentService := service.NewCommentService(commentRepo)
	visitorService := service.NewVisitorService(visitorRepo, commentRepo)
	postLikeService := service.NewPostLikeService(postLikeRepo, postRepo)
	devlogService := service.NewDevlogService(devlogRepo)
	fileService := service.NewLocalFileService()

	// 7. 创建Gin引擎
	engine := gin.Default()

	// 配置路由
	router.Setup(engine, postService, columnService, commentService, visitorService, postLikeService, fileService, devlogService)

	// 8. 启动服务
	addr := cfg.Address()
	log.Printf("服务器启动成功: http://%s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}

	// 防止编译提前退出
	log.Println("按 Ctrl+C 停止服务")
}