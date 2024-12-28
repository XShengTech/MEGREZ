package admin

import (
	"megrez/routers/api/v1/admin/images"
	instances "megrez/routers/api/v1/admin/instance"
	"megrez/routers/api/v1/admin/servers"
	"megrez/routers/api/v1/admin/users"
	"megrez/routers/api/v1/middleware"

	"github.com/kataras/iris/v12/core/router"
)

func InitAdmin(party router.Party) {
	party.Use(middleware.Auth, middleware.AuthCheck, middleware.AdminCheck)

	users.InitUser(party.Party("/users"))
	servers.InitServer(party.Party("/servers"))
	instances.InitInstances(party.Party("/instances"))
	images.InitImages(party.Party("/images"))
}
