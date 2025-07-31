package routes

import (
	"fmt"
	"news-release/internal/config"
	"news-release/internal/database"
	"news-release/internal/middleware"

	articlectr "news-release/internal/article/controller"
	articlerepo "news-release/internal/article/repository"
	articlesvc "news-release/internal/article/service"

	noticectr "news-release/internal/notice/controller"
	noticerepo "news-release/internal/notice/repository"
	noticesvc "news-release/internal/notice/service"

	userctr "news-release/internal/user/controller"
	userrepo "news-release/internal/user/repository"
	usersvc "news-release/internal/user/service"

	filectr "news-release/internal/file/controller"
	filerepo "news-release/internal/file/repository"
	filesvc "news-release/internal/file/service"

	msgctr "news-release/internal/message/controller"
	msgrepo "news-release/internal/message/repository"
	msgsvc "news-release/internal/message/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRoutes 注册中间件和路由
func SetupRoutes(cfg *config.Config, router *gin.Engine) {
	// 注册中间件
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// 连接数据库
	DSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)
	db, err := database.NewDatabase(DSN)
	if err != nil {
		logrus.Panic(err)
	}

	// 创建MinIO存储实例
	minioRepo, err := filerepo.NewMinIORepository(
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKeyID,
		cfg.MinIO.SecretAccessKey,
		cfg.MinIO.UseSSL,
		cfg.MinIO.BucketName,
	)
	if err != nil {
		logrus.Panic("创建MinIO存储实例失败: ", err)
	}

	// 初始化依赖
	// 初始化仓库
	articleRepo := articlerepo.NewArticleRepository(db)
	fieldTypeRepo := articlerepo.NewFieldTypeRepository(db)
	noticeRepo := noticerepo.NewNoticeRepository(db)
	fileRepo := filerepo.NewFileRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	msgRepo := msgrepo.NewMessageRepository(db)

	// 初始化服务
	articleService := articlesvc.NewArticleService(articleRepo)
	fieldService := articlesvc.NewFieldTypeService(fieldTypeRepo)
	noticeService := noticesvc.NewNoticeService(noticeRepo)
	fileService := filesvc.NewFileService(minioRepo, fileRepo)
	userService := usersvc.NewUserService(userRepo, cfg)
	msgService := msgsvc.NewMessageService(msgRepo)

	// 初始化控制器
	articleController := articlectr.NewArticleController(articleService)
	fieldTypeController := articlectr.NewFieldTypeController(fieldService)
	noticeController := noticectr.NewNoticeController(noticeService)
	fileController := filectr.NewFileController(fileService)
	userController := userctr.NewUserController(userService)
	msgController := msgctr.NewMessageController(msgService)

	// API分组
	api := router.Group("/api")
	{
		// articles
		articles := api.Group("/articles")
		{
			articles.GET("", articleController.ListArticle)
			articles.GET("/:articleID", articleController.GetArticleContent)
		}
		// 领域类型相关路由
		policyFieldType := api.Group("/fieldType")
		{
			policyFieldType.GET("", fieldTypeController.GetFieldType)
		}
		// 公告相关路由
		notice := api.Group("/notice")
		{
			notice.GET("", noticeController.ListNotice)
			notice.GET("/:ID", noticeController.GetNoticeContent)
		}
		// 用户相关路由
		user := api.Group("/user")
		{
			user.POST("/login", userController.Login)
		}
		// 文件上传路由
		file := api.Group("/file")
		// file.Use(middleware.AuthMiddleware(cfg))
		{
			// 上传文件需进行身份验证
			file.POST("/upload", middleware.AuthMiddleware(cfg), fileController.UploadFile)
		}
		// 消息相关路由
		message := api.Group("/message")
		message.Use(middleware.AuthMiddleware(cfg))
		{
			message.GET("", msgController.ListMessage)
			message.GET("/:messageID", msgController.GetMessageContent)
		}
	}
}
