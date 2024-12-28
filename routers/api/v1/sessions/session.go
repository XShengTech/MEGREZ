package sessions

import (
	"time"

	"github.com/kataras/iris/v12/sessions"
)

func Session() *sessions.Sessions {
	return sessions.New(sessions.Config{
		Cookie:  "session_id",
		Expires: 7 * 24 * time.Hour, // 7 Days
	})
}
