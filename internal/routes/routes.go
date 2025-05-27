package routes

import (
	"news-release/internal/config"
	"news-release/internal/controller"
	"news-release/internal/middleware"
	"news-release/internal/repository"
	"news-release/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 注册路由
func SetupRoutes(cfg *config.Config, router *gin.Engine) {
	// 连接数据库
	db, err := repository.NewDatabase(cfg.DBURL)
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	// 初始化依赖
	// 初始化仓库
	exampleRepo := repository.NewExampleRepository(db)
	// 初始化服务
	exampleService := service.NewExampleService(exampleRepo)
	// 初始化控制器
	exampleController := controller.NewExampleController(exampleService)

	// 注册中间件
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// API分组
	api := router.Group("/api")
	{
		// example仅用于示例及测试
		example := api.Group("/example")
		{
			example.GET("/ListExample", exampleController.ListExample)
		}
		{
			// 其他API路由...
		}
	}
}
