package user

import (
	"megrez/routers/api/v1/middleware"
	"megrez/routers/api/v1/user/images"
	"megrez/routers/api/v1/user/instances"
	"megrez/routers/api/v1/user/servers"
	"megrez/services/logger"

	_logger "megrez/libs/logger"

	"github.com/kataras/iris/v12/core/router"
)

var l *_logger.LoggerStruct

func InitUser(party router.Party) {
	l = logger.Logger.Clone()
	l.SetModel("Http.API.V1.User")

	party.Use(middleware.Auth)

	party.Post("/login", loginHandler)
	party.Get("/logout", middleware.AuthCheck, logoutHandler)
	party.Post("/register", registerHandler)
	party.Get("/profile", middleware.AuthCheck, profileHandler)
	party.Post("/password", middleware.AuthCheck, resetPasswordHandler)
	party.Post("/email", middleware.AuthCheck, resetEmailHandler)
	party.Get("/verify/{code:string}", verifyHandler)
	party.Post("/verify", middleware.AuthCheck, verifySendHandler)

	servers.InitServers(party.Party("/servers"))
	instances.InitInstances(party.Party("/instances"))
	images.InitImages(party.Party("/images"))
}
