package http

import (
	"context"
	"fmt"
	_logger "megrez/libs/logger"
	"megrez/routers"
	"megrez/services/config"
	"megrez/services/logger"
	"time"

	"github.com/kataras/iris/v12"
)

var l *_logger.LoggerStruct

func Start() {
	l = logger.Logger.Clone()
	l.SetModel("HTTP")
	l.SetFunction("Start")

	hc := config.GetHttpAddress()

	app := iris.New()
	app.Use(iris.Compression)

	idleConnsClosed := make(chan struct{})
	iris.RegisterOnInterrupt(func() {
		fmt.Println()
		l.SetFunction("CloseHttp")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := app.Shutdown(ctx)
		if err != nil {
			l.Error("Http Shutdown err, Error: %v", err.Error())
		}
		l.Info("Http Closed")

		close(idleConnsClosed)
	})

	routers.InitRouter(app)

	l.Info("Http Server Listening on http://%s", hc)
	l.Info("Press CTRL+C to shut down.")
	err := app.Listen(
		hc,
		iris.WithConfiguration(iris.Configuration{
			DisableStartupLog:   true,
			LogLevel:            "disable",
			Charset:             "UTF-8",
			EnableOptimizations: true,
			RemoteAddrHeaders: []string{
				"X-Real-Ip",
				"X-Forwarded-For",
				"CF-Connecting-IP",
				"True-Client-Ip",
				"X-Appengine-Remote-Addr",
			},
		}),
		iris.WithoutServerError(iris.ErrServerClosed),
	)
	if err != nil {
		l.Error("Http Server Listen err, Error: %v", err.Error())
		return
	}
	<-idleConnsClosed
}
