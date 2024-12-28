package servers

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"

	"github.com/kataras/iris/v12"
)

func modifyHandler(ctx iris.Context) {
	l.SetFunction("modifyHandler")

	id, err := ctx.Params().GetUint("id")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	var req serverStruct
	if err := ctx.ReadJSON(&req); err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	server := models.Servers{
		ID: id,
	}
	result := database.DB.First(&server)
	if result.Error != nil {
		l.Error("query server error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeAdminServerAddEditError, iris.StatusInternalServerError)
		return
	}

	if req.Name != "" && req.Name != server.Name {
		server.Name = req.Name
	}
	if req.IP != "" && req.IP != server.IP {
		server.IP = req.IP
	}
	if req.Port != 0 && req.Port != server.Port {
		server.Port = req.Port
	}
	if req.Apikey != "" && req.Apikey != server.Apikey {
		server.Apikey = req.Apikey
	}
	if req.GpuType != "" && req.GpuType != server.GpuType {
		server.GpuType = req.GpuType
	}
	if req.GpuNum != 0 && req.GpuNum != server.GpuNum {
		server.GpuNum = req.GpuNum
	}
	if req.GpuDirverVersion != "" && req.GpuDirverVersion != server.GpuDriverVersion {
		server.GpuDriverVersion = req.GpuDirverVersion
	}
	if req.GpuCudaVersion != "" && req.GpuCudaVersion != server.GpuCudaVersion {
		server.GpuCudaVersion = req.GpuCudaVersion
	}
	if req.CpuCpuntPerGpu != 0 && req.CpuCpuntPerGpu != server.CpuCountPerGpu {
		server.CpuCountPerGpu = req.CpuCpuntPerGpu
	}
	if req.MemoryPerGpu != 0 && req.MemoryPerGpu != server.MemoryPerGpu {
		server.MemoryPerGpu = req.MemoryPerGpu
	}
	if req.VolumeTotal != 0 && req.VolumeTotal != server.VolumeTotal {
		server.VolumeTotal = req.VolumeTotal
	}
	if req.Price != 0 && req.Price != server.Price {
		server.Price = req.Price
	}
	if req.PriceVolume != 0 && req.PriceVolume != server.PriceVolume {
		server.PriceVolume = req.PriceVolume
	}

	result = database.DB.Save(&server)
	if result.Error != nil {
		l.Error("save server error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeAdminServerAddEditError, iris.StatusInternalServerError)
		return
	}

	middleware.Success(ctx)
}
