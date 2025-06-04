package utils

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// 初始化日志记录器,使用 logrus 作为日志记录器
func InitLogger() {
	// 创建日志目录
	os.MkdirAll("../logs/info", 0755)
	os.MkdirAll("../logs/error", 0755)

	infoLogPath := "../logs/info/info.log"
	infoWriter, err := rotatelogs.New(
		infoLogPath+".%Y%m%d",
		rotatelogs.WithLinkName(infoLogPath),
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 保留7天
		rotatelogs.WithRotationTime(24*time.Hour), // 每天分割
	)
	if err != nil {
		logrus.Fatalf("配置 Info 日志分割器失败: %v", err)
	}

	// 创建 Error 级别日志的分割器（处理 Error、Fatal 和 Panic 级别）
	errLogPath := "../logs/error/err.log"
	errorWriter, err := rotatelogs.New(
		errLogPath+".%Y%m%d",
		rotatelogs.WithLinkName(errLogPath),
		rotatelogs.WithMaxAge(30*24*time.Hour),    // 保留30天
		rotatelogs.WithRotationTime(24*time.Hour), // 每天分割
	)
	if err != nil {
		logrus.Fatalf("配置 Error 日志分割器失败: %v", err)
	}

	// 添加日志钩子
	logrus.AddHook(lfshook.NewHook(
		// 定义不同日志级别对应的输出位置
		lfshook.WriterMap{
			logrus.InfoLevel:  infoWriter,
			logrus.WarnLevel:  infoWriter,
			logrus.ErrorLevel: errorWriter,
			logrus.FatalLevel: errorWriter,
			logrus.PanicLevel: errorWriter,
		},
		// 设置日志格式为 JSON 格式，并指定时间戳的格式为 RFC3339 标准
		&logrus.JSONFormatter{TimestampFormat: time.RFC3339},
	))
	// 设置 logrus 为默认日志记录器
	log.SetOutput(logrus.StandardLogger().Writer())
	gin.DefaultWriter = logrus.StandardLogger().Writer()
}
