package v1

import (
	"megrez/routers/api/v1/admin"
	"megrez/routers/api/v1/user"

	"github.com/kataras/iris/v12/core/router"
)

func InitApiV1(party router.Party) {
	admin.InitAdmin(party.Party("/admin"))
	user.InitUser(party.Party("/user"))
}
