package users

import (
	_logger "megrez/libs/logger"
	"megrez/routers/api/v1/middleware"
	"megrez/services/logger"

	"github.com/kataras/iris/v12/core/router"
)

var l *_logger.LoggerStruct

func InitUser(party router.Party) {
	l = logger.Logger.Clone()
	l.SetModel("Http.API.V1.Admin.Users")

	party.Get("/", listHandler)
	party.Get("/{id:uint}", detailHandler)
	// party.Post("/", addHandler)
	party.Post("/{id:uint}", middleware.SuperAdminCheck, modifyHandler)
	party.Delete("/{id:uint}", middleware.SuperAdminCheck, deleteHandler)
}
