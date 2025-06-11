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

	// 初始化依赖
	// 初始化仓库
	exampleRepo := repository.NewExampleRepository(db)
	policyRepo := repository.NewPolicyRepository(db)
	// 初始化服务
	exampleService := service.NewExampleService(exampleRepo)
	policyService := service.NewPolicyService(policyRepo)
	// 初始化控制器
	exampleController := controller.NewExampleController(exampleService)
	policyController := controller.NewPolicyController(policyService)

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
	}
}
