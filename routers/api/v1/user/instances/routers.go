package instances

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/logger"

	_logger "megrez/libs/logger"

	"github.com/kataras/iris/v12/core/router"
)

var l *_logger.LoggerStruct

type instanceStruct struct {
	models.Instances
	ServerName     string `json:"server_name"`
	GpuType        string `json:"gpu_type"`
	GpuNum         int    `json:"gpu_num"`
	GpuUsed        int    `json:"gpu_used"`
	CpuCountPerGpu int    `json:"cpu_count_per_gpu"`
	MemoryPerGpu   int    `json:"memory_per_gpu"`
}

func InitInstances(party router.Party) {
	l = logger.Logger.Clone()
	l.SetModel("Http.API.V1.User.Instance")

	party.Use(middleware.AuthCheck, middleware.UserCheck)

	party.Get("/", listHandler)
	party.Get("/{id:uint}", detailHandler)
	party.Put("/{id:uint}", controlHandler)
	party.Post("/", addHandler)
	party.Post("/{id:uint}", modifyHandler)
	party.Post("/{id:uint}/label", labelHandler)
	party.Delete("/{id:uint}", deleteHandler)
}
