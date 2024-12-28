package routers

import (
	_logger "megrez/libs/logger"
	v1 "megrez/routers/api/v1"
	"megrez/routers/index"
	"megrez/services/logger"

	"github.com/kataras/iris/v12"
)

var l *_logger.LoggerStruct

func InitRouter(app *iris.Application) {
	l = logger.Logger.Clone()
	l.SetModel("HTTP")

	app.Use(middleware)
	v1.InitApiV1(app.Party("/api/v1"))
	index.InitIndex(app)
}
