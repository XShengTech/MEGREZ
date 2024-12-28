package logger

import "megrez/libs/logger"

var Logger *logger.LoggerStruct

func init() {
	Logger, _ = logger.NewLogger(logger.DEBUG, "stdout")
}

func InitLogger(level, path string) {
	Logger, _ = logger.NewLogger(level, path)
}
