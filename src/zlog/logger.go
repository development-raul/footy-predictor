package zlog

import (
	"go.uber.org/zap"
	"os"
	"strings"
)

var Logger *zap.SugaredLogger

func init() {
	appEnv := strings.ToLower(os.Getenv("APP_ENV"))
	logPath := os.Getenv("ERROR_LOG_PATH")
	if logPath == "" {
		logPath = "application.log"
	}

	config := zap.Config{}

	if appEnv == "prod" || appEnv == "production" {
		config = zap.NewProductionConfig()
		config.OutputPaths = []string{"stderr", logPath}
	} else {
		config = zap.NewDevelopmentConfig()
		config.OutputPaths = []string{"stderr", logPath}
	}

	standardLogger, err := config.Build()
	if err != nil {
		panic(err)
	}

	Logger = standardLogger.Sugar() // Use sugar logger
}
