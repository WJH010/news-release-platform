package main

import (
	"net/http"
	"news-release/internal/config"
	"news-release/internal/routes"
	"news-release/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化日志记录器
	utils.InitLogger()

	// 加载配置
	cfg, err := config.LoadConfig("../config.yaml")
	if err != nil {
		logrus.Fatalf("服务器启动失败: %v", err)
	}

	// 设置Gin模式
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建默认的Gin引擎，但不使用默认中间件
	router := gin.New()

	// 初始化依赖及注册路由
	routes.SetupRoutes(cfg, router)

	PORT := cfg.App.Port
	//log.Printf("服务器运行在端口 %d", PORT)
	logrus.Infof("服务器运行在端口 %d", PORT)
	if err := http.ListenAndServe(":"+strconv.Itoa(PORT), router); err != nil {
		//log.Fatalf("服务器启动失败: %v", err)
		logrus.Fatalf("服务器启动失败: %v", err)
	}
}
