package images

import (
	"megrez/routers/api/v1/middleware"
	"megrez/services/logger"

	"github.com/kataras/iris/v12/core/router"

	_logger "megrez/libs/logger"
)

const imagesKey = "images"

var l *_logger.LoggerStruct

func InitImages(party router.Party) {
	l = logger.Logger.Clone()
	l.SetModel("Http.API.V1.Admin.Images")

	party.Get("/", listHandler)
	party.Post("/", middleware.SuperAdminCheck, modifyHandler)
}
