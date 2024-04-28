package fastgo

import (
	"go.uber.org/zap"
	"log"
)

// 获取log记录器
var logger *zap.Logger

func Logger() *zap.Logger {
	return logger
}

func init() {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	logger = l
	zap.ReplaceGlobals(logger)
}

func LoggerSync() {
	_ = logger.Sync()
}

// serverErrorLogger
func serverErrorLogger() *log.Logger {
	logAt, err := zap.NewStdLogAt(logger, zap.ErrorLevel)
	if err != nil {
		return nil
	}
	return logAt
}
