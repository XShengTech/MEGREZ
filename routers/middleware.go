package routers

import (
	"github.com/kataras/iris/v12"
)

func middleware(ctx iris.Context) {
	ctx.Next()

	l.Info("IP: %s - %s %d %s %d - User-Agent: %s", ctx.RemoteAddr(), ctx.Request().Method, ctx.GetStatusCode(), ctx.Request().RequestURI, ctx.Request().ContentLength, ctx.Request().UserAgent())
}
