package routes

import (
	"fmt"
	"news-release/internal/config"
	"news-release/internal/controller"
	"news-release/internal/middleware"
	"news-release/internal/repository"
	"news-release/internal/service"

	userctr "news-release/internal/controller/user"
	userrepo "news-release/internal/repository/user"
	usersvc "news-release/internal/service/user"

	msgctr "news-release/internal/controller/message"
	msgrepo "news-release/internal/repository/message"
	msgsvc "news-release/internal/service/message"

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
	db, err := repository.NewDatabase(DSN)
	if err != nil {
		logrus.Panic("数据库连接失败: ", err)
	}

	// 创建MinIO存储实例
	minioRepo, err := repository.NewMinIORepository(
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
	exampleRepo := repository.NewExampleRepository(db)
	articleRepo := repository.NewArticleRepository(db)
	fieldTypeRepo := repository.NewFieldTypeRepository(db)
	noticeRepo := repository.NewNoticeRepository(db)
	fileRepo := repository.NewFileRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	adminRepo := repository.NewAdminUserRepository(db)
	msgRepo := msgrepo.NewMessageRepository(db)

	// 初始化服务
	exampleService := service.NewExampleService(exampleRepo)
	articleService := service.NewArticleService(articleRepo)
	fieldService := service.NewFieldTypeService(fieldTypeRepo)
	noticeService := service.NewNoticeService(noticeRepo)
	fileService := service.NewFileService(minioRepo, fileRepo)
	userService := usersvc.NewUserService(userRepo, cfg)
	adminService := service.NewAdminUserService(adminRepo)
	msgService := msgsvc.NewMessageService(msgRepo)

	// 初始化控制器
	exampleController := controller.NewExampleController(exampleService)
	articleController := controller.NewArticleController(articleService)
	fieldTypeController := controller.NewFieldTypeController(fieldService)
	noticeController := controller.NewNoticeController(noticeService)
	fileController := controller.NewFileController(fileService)
	userController := userctr.NewUserController(userService)
	adminController := controller.NewAdminController(adminService, cfg)
	msgController := msgctr.NewMessageController(msgService)

	// API分组
	api := router.Group("/api")
	{
		// example仅用于示例及测试
		example := api.Group("/example")
		{
			example.GET("", exampleController.ListExample)
		}
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
		}
		// 用户相关路由
		user := api.Group("/user")
		{
			user.POST("/login", userController.Login)
		}
		admin := api.Group("/admin")
		{
			admin.POST("/login", adminController.AdminLogin) // 管理系统登录接口
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
		}
	}
}
