package routes

import (
	"fmt"
	"news-release/internal/config"
	"news-release/internal/controller"
	"news-release/internal/middleware"
	"news-release/internal/repository"
	"news-release/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRoutes 注册路由
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
	policyRepo := repository.NewPolicyRepository(db)
	newsRope := repository.NewNewsRepository(db)
	fieldTypeRepo := repository.NewFieldTypeRepository(db)
	noticeRepo := repository.NewNoticeRepository(db)
	fileRepo := repository.NewFileRepository(db)
	userRepo := repository.NewUserRepository(db)
	adminRepo := repository.NewAdminUserRepository(db)

	// 初始化服务
	exampleService := service.NewExampleService(exampleRepo)
	policyService := service.NewPolicyService(policyRepo)
	fieldService := service.NewFieldTypeService(fieldTypeRepo)
	noticeService := service.NewNoticeService(noticeRepo)
	newsService := service.NewNewsService(newsRope)
	fileService := service.NewFileService(minioRepo, fileRepo)
	userService := service.NewUserService(userRepo, cfg)
	adminService := service.NewAdminUserService(adminRepo)

	// 初始化控制器
	exampleController := controller.NewExampleController(exampleService)
	policyController := controller.NewPolicyController(policyService)
	newsController := controller.NewNewsController(newsService)
	fieldTypeController := controller.NewFieldTypeController(fieldService)
	noticeController := controller.NewNoticeController(noticeService)
	fileController := controller.NewFileController(fileService)
	userController := controller.NewUserController(userService)
	adminController := controller.NewAdminController(adminService, cfg)

	// API分组
	api := router.Group("/api")
	{
		// example仅用于示例及测试
		example := api.Group("/example")
		{
			example.GET("/ListExample", exampleController.ListExample)
		}
		// policy相关路由
		policy := api.Group("/policy")
		{
			policy.GET("/ListPolicy", policyController.ListPolicy)
			policy.GET("/GetPolicyContent/:id", policyController.GetPolicyContent)
		}
		// news相关路由
		news := api.Group("/news")
		{
			news.GET("/ListNews", newsController.GetNewsList)
			news.GET("/GetNewsContent/:id", newsController.GetNewsContent)
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
		// 文件上传路由
		file := api.Group("/file")
		{
			file.POST("/upload", fileController.UploadFile)
		}
		// 用户相关路由
		user := api.Group("/user")
		{
			user.GET("/login", userController.Login)
		}
		admin := api.Group("/admin")
		{
			admin.POST("/login", adminController.AdminLogin) // 管理系统登录接口
		}
	}
}
