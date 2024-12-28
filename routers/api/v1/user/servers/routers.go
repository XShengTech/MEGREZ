package servers

import (
	"megrez/services/logger"

	"github.com/kataras/iris/v12/core/router"

	_logger "megrez/libs/logger"
)

var l *_logger.LoggerStruct

func InitServers(party router.Party) {
	l = logger.Logger.Clone()
	l.SetModel("Http.API.V1.User.Servers")

	party.Get("/", listHandler)
	party.Get("/{id:uint}", detailHandler)
}
