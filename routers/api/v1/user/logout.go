package user

import (
	"megrez/routers/api/v1/middleware"
	"megrez/routers/api/v1/sessions"
	"megrez/services/redis"
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
)

func logoutHandler(ctx iris.Context) {
	sess := sessions.Session()
	session := sess.Start(ctx)
	sess.UseDatabase(redis.DB)
	sessionId := session.ID()

	session.Destroy()
	ctx.SetCookie(&http.Cookie{
		Name:    "session_id",
		Value:   sessionId,
		Expires: time.Now(),
	})
	middleware.Success(ctx)
}
