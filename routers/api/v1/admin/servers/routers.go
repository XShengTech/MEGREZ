package servers

import (
	"github.com/kataras/iris/v12/core/router"

	"megrez/routers/api/v1/middleware"
	"megrez/services/logger"

	_logger "megrez/libs/logger"
)

var l *_logger.LoggerStruct

type serverStruct struct {
	Name             string `json:"name,omitempty"`
	IP               string `json:"ip,omitempty"`
	Port             int    `json:"port,omitempty"`
	Apikey           string `json:"apikey,omitempty"`
	GpuType          string `json:"gpu_type,omitempty"`
	GpuNum           int    `json:"gpu_num,omitempty"`
	GpuDirverVersion string `json:"gpu_driver_version,omitempty"`
	GpuCudaVersion   string `json:"gpu_cuda_version,omitempty"`

	CpuCpuntPerGpu int     `json:"cpu_count_per_gpu,omitempty"`
	MemoryPerGpu   int     `json:"memory_per_gpu,omitempty"` // Unit `GB`
	VolumeTotal    int     `json:"volume_total,omitempty"`   // Unit `GB`
	Price          float64 `json:"price,omitempty"`          // 1 GPU Per Hour
	PriceVolume    float64 `json:"price_volume,omitempty"`   // 1GB Per Hour
}

func InitServer(party router.Party) {
	l = logger.Logger.Clone()
	l.SetModel("Http.API.V1.Admin.Servers")

	party.Get("/", listHandler)
	party.Get("/{id:uint}", detailHandler)
	party.Post("/", middleware.SuperAdminCheck, addHandler)
	party.Post("/{id:uint}", middleware.SuperAdminCheck, modifyHandler)
	party.Delete("/{id:uint}", middleware.SuperAdminCheck, deleteHandler)
}
