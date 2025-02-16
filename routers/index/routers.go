package index

import (
	"github.com/kataras/iris/v12"
)

func InitIndex(app *iris.Application) {
	app.Use(cors)
	app.HandleDir("/", GetWebFS(), iris.DirOptions{
		IndexName: "/index.html",
		Compress:  true,
		SPA:       true,
	})
}
