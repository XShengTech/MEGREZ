package redis

import (
	"strconv"
	"time"

	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
	goRedis "github.com/redis/go-redis/v9"

	_logger "megrez/libs/logger"
	"megrez/services/config"
	"megrez/services/logger"
)

var l *_logger.LoggerStruct
var DB *redis.Database
var RawDB *goRedis.Client

func Connect() (err error) {
	l = logger.Logger.Clone()
	l.SetModel("redis")
	l.SetFunction("Connect")

	rd := config.GetRedis()
	DB = redis.New(redis.Config{
		Network:   "tcp",
		Addr:      rd.Host + ":" + strconv.Itoa(rd.Port),
		Timeout:   time.Duration(30) * time.Second,
		MaxActive: 10,
		Password:  rd.Password,
		Database:  strconv.Itoa(rd.Database),
		Prefix:    "",
	})

	RawDB = goRedis.NewClient(&goRedis.Options{
		Addr:     rd.Host + ":" + strconv.Itoa(rd.Port),
		Password: rd.Password,
		DB:       rd.Database,
	})

	l.Info("Redis connected")

	return nil
}

func Close() (err error) {
	l.SetFunction("Close")

	if err = DB.Close(); err != nil {
		l.Error("Failed to close Redis Client, Error: %v", err)
		return
	}

	if err = RawDB.Close(); err != nil {
		l.Error("Failed to close Redis Client, Error: %v", err)
		return
	}

	l.Info("Redis closed")
	return
}
