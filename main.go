package main

import (
	"flag"
	"megrez/services/config"
	"megrez/services/database"
	"megrez/services/dispatcher"
	"megrez/services/http"
	"megrez/services/instanceController"
	"megrez/services/logger"
	"megrez/services/redis"
	"megrez/services/system"
	"runtime"
)

var (
	BRANCH    string
	VERSION   string
	COMMIT    string
	GoVersion string
	BuildTime string
)

var (
	configFilePath = "config.yml"
	l              = logger.Logger.Clone()
)

func main() {
	l.Info("Branch: %s", BRANCH)
	l.Info("Version: %s", VERSION)
	l.Info("Commit: %s", COMMIT)
	l.Info("Go Version: %s", GoVersion)
	l.Info("Build Time: %s", BuildTime)

	flag.StringVar(&configFilePath, "c", "config.yml", "config file path")
	flag.Parse()

	config.InitConfig(configFilePath)

	defer func() {
		redis.Close()

		l.Close()
	}()
	defer func() {
		defer func() {
			if err := recover(); err != nil {
				l.Error("Panic: %v", err)
				buf := make([]byte, 1024)
				n := runtime.Stack(buf, false)
				l.Error("Stack trace: \n%s", buf[:n])
			}
		}()
	}()

	logger.InitLogger(config.GetLogLevel(), config.GetLogFile())

	database.Connect()
	redis.Connect()
	instanceController.InitInstanceController()
	dispatcher.Init()

	go system.Check()

	http.Start()
}
