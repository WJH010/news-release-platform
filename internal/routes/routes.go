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

	eventctr "news-release/internal/event/controller"
	eventrepo "news-release/internal/event/repository"
	eventsvc "news-release/internal/event/service"

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
	industryRepo := userrepo.NewIndustryRepository(db)
	msgRepo := msgrepo.NewMessageRepository(db)
	msgType := msgrepo.NewMessageTypeRepository(db)
	eventRepo := eventrepo.NewEventRepository(db)

	// 初始化服务
	articleService := articlesvc.NewArticleService(articleRepo, fileRepo)
	fieldService := articlesvc.NewFieldTypeService(fieldTypeRepo)
	noticeService := noticesvc.NewNoticeService(noticeRepo)
	fileService := filesvc.NewFileService(minioRepo, fileRepo)
	userService := usersvc.NewUserService(userRepo, cfg)
	industryService := usersvc.NewIndustryService(industryRepo)
	msgService := msgsvc.NewMessageService(msgRepo)
	msgTypeService := msgsvc.NewMessageTypeService(msgType)
	eventService := eventsvc.NewEventService(eventRepo, userRepo, fileRepo)

	// 初始化控制器
	articleController := articlectr.NewArticleController(articleService)
	fieldTypeController := articlectr.NewFieldTypeController(fieldService)
	noticeController := noticectr.NewNoticeController(noticeService)
	fileController := filectr.NewFileController(fileService)
	userController := userctr.NewUserController(userService)
	industryController := userctr.NewIndustryController(industryService)
	msgController := msgctr.NewMessageController(msgService)
	msgTypeController := msgctr.NewMessageTypeController(msgTypeService)
	eventController := eventctr.NewEventController(eventService)

	// API分组
	api := router.Group("/api")
	{
		// articles
		articles := api.Group("/articles")
		{
			// 公开接口 - 无需认证
			articles.GET("", articleController.ListArticle)
			articles.GET("/:id", articleController.GetArticleContent)
			// 需要认证的用户接口
			authArticles := articles.Group("")
			authArticles.Use(middleware.AuthMiddleware(cfg))
			{
				// 管理员接口 - 在认证基础上增加角色校验
				adminArticles := authArticles.Group("")
				adminArticles.Use(middleware.RoleMiddleware(middleware.RoleAdmin))
				{
					adminArticles.POST("/create", articleController.CreateArticle)
					adminArticles.PUT("/update/:id", articleController.UpdateArticle)
					adminArticles.DELETE("/delete/:id", articleController.DeleteArticle)
				}
			}
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
			notice.GET("/:id", noticeController.GetNoticeContent)
		}
		// 用户相关路由
		user := api.Group("/user")
		{
			// 公开接口 - 无需认证
			user.POST("/login", userController.Login)
			// 需要认证的用户接口
			authUser := user.Group("")
			authUser.Use(middleware.AuthMiddleware(cfg))
			{
				authUser.PUT("/update", middleware.AuthMiddleware(cfg), userController.UpdateUserInfo)
				authUser.GET("/info", middleware.AuthMiddleware(cfg), userController.GetUserInfo)
				// 管理员接口 - 在认证基础上增加角色校验
				adminUser := authUser.Group("")
				adminUser.Use(middleware.RoleMiddleware(middleware.RoleAdmin))
				{

				}
			}
		}
		// 行业路由
		industry := api.Group("/industry")
		{
			// 公开接口 - 无需认证
			industry.GET("", industryController.ListIndustries)
			// 需要认证的用户接口
			authIndustry := industry.Group("")
			authIndustry.Use(middleware.AuthMiddleware(cfg))
			{
				// 管理员接口 - 在认证基础上增加角色校验
				adminIndustry := authIndustry.Group("")
				adminIndustry.Use(middleware.RoleMiddleware(middleware.RoleAdmin))
				{
					adminIndustry.POST("/create", industryController.CreateIndustry)
					adminIndustry.PUT("/update/:id", industryController.UpdateIndustry)
				}
			}
		}
		// 文件上传路由
		file := api.Group("/file")
		file.Use(middleware.AuthMiddleware(cfg))
		{
			file.POST("/upload", fileController.UploadFile)
			file.DELETE("/deleteImage/:id", fileController.DeleteImage)

		}
		// 消息相关路由
		message := api.Group("/message")
		message.Use(middleware.AuthMiddleware(cfg))
		{
			message.GET("/:id", msgController.GetMessageContent)
			message.GET("/unreadMessageCount", msgController.GetUnreadMessageCount)
			message.PUT("/markAllAsRead", msgController.MarkAllMessagesAsRead)
			message.GET("/userMessageGroups", msgController.ListUserMessageGroups)
			message.GET("/byGroups", msgController.ListMsgByGroups)
		}
		messageType := api.Group("/messageType")
		{
			messageType.GET("", msgTypeController.ListMessageType)
		}
		// 活动相关路由
		event := api.Group("/event")
		{
			// 公开接口 - 无需认证
			event.GET("", eventController.ListEvent)
			event.GET("/:id", eventController.GetEventDetail)

			// 需要认证的用户接口
			authEvent := event.Group("")
			authEvent.Use(middleware.AuthMiddleware(cfg))
			{
				authEvent.POST("/registration", eventController.RegistrationEvent)
				authEvent.GET("/isUserRegistered/:id", eventController.IsUserRegistered)
				authEvent.DELETE("/cancelRegistration/:id", eventController.CancelRegistrationEvent)
				authEvent.GET("/userRegisteredEvents", eventController.ListUserRegisteredEvents)

				// 管理员接口 - 在认证基础上增加角色校验
				adminEvent := authEvent.Group("")
				adminEvent.Use(middleware.RoleMiddleware(middleware.RoleAdmin))
				{
					adminEvent.POST("/create", eventController.CreateEvent)
					adminEvent.PUT("/update/:id", eventController.UpdateEvent)
					adminEvent.DELETE("/delete/:id", eventController.DeleteEvent)
					adminEvent.GET("/regUsers/:id", eventController.ListEventRegisteredUsers)
				}
			}
		}
	}
}
