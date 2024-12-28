package middleware

import (
	"megrez/routers/api/v1/sessions"
	"megrez/services/redis"

	"github.com/kataras/iris/v12"
)

func Auth(ctx iris.Context) {
	sess := sessions.Session()
	sess.UseDatabase(redis.DB)
	session := sess.Start(ctx)

	auth, _ := session.GetBoolean("authenticated")
	if !auth {
		ctx.Values().Set("authenticated", false)
		ctx.Next()
		return
	} else {
		ctx.Values().Set("authenticated", true)
	}

	userId, _ := session.GetInt("userId")
	ctx.Values().Set("userId", userId)
	ctx.Values().Set("email", session.GetString("email"))

	role, _ := session.GetInt("role")
	ctx.Values().Set("role", role)

	ctx.Next()
}

func AuthCheck(ctx iris.Context) {
	if ctx.Values().Get("authenticated") == nil || !ctx.Values().Get("authenticated").(bool) {
		Error(ctx, CodeUnauthorized, iris.StatusUnauthorized)
		return
	}

	ctx.Next()
}

func UserCheck(ctx iris.Context) {
	usertype, _ := ctx.Values().GetInt("role")
	if usertype < 1 {
		Error(ctx, CodeForbidden, iris.StatusForbidden)
		return
	}

	ctx.Next()
}

func AdminCheck(ctx iris.Context) {
	usertype, _ := ctx.Values().GetInt("role")
	if usertype < 2 {
		Error(ctx, CodeForbidden, iris.StatusForbidden)
		return
	}

	ctx.Next()
}

func SuperAdminCheck(ctx iris.Context) {
	usertype, _ := ctx.Values().GetInt("role")
	if usertype < 3 {
		Error(ctx, CodeForbidden, iris.StatusForbidden)
		return
	}

	ctx.Next()
}
