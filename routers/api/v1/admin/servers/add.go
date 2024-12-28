package servers

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"
	"megrez/services/redis"
	"strconv"

	"github.com/kataras/iris/v12"
)

func addHandler(ctx iris.Context) {
	l.SetFunction("addHandler")

	var s serverStruct
	if err := ctx.ReadJSON(&s); err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	if s.Name == "" || s.IP == "" || s.Port == 0 || s.Apikey == "" || s.GpuType == "" || s.GpuNum == 0 || s.GpuDirverVersion == "" || s.GpuCudaVersion == "" || s.CpuCpuntPerGpu == 0 || s.MemoryPerGpu == 0 || s.VolumeTotal == 0 {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	l.Debug("add server: %+v", s)

	server := models.Servers{
		Name:             s.Name,
		IP:               s.IP,
		Port:             s.Port,
		Apikey:           s.Apikey,
		GpuType:          s.GpuType,
		GpuNum:           s.GpuNum,
		GpuDriverVersion: s.GpuDirverVersion,
		GpuCudaVersion:   s.GpuCudaVersion,

		CpuCountPerGpu: s.CpuCpuntPerGpu,
		MemoryPerGpu:   s.MemoryPerGpu,
		VolumeTotal:    s.VolumeTotal,
		Price:          s.Price,
		PriceVolume:    s.PriceVolume,
	}
	result := database.DB.Create(&server)
	if result.Error != nil {
		l.Error("add server error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeAdminServerAddEditError, iris.StatusInternalServerError)
		return
	}

	redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(server.ID)), int64(server.GpuNum))
	redis.RawDB.IncrBy(ctx, "remain_volume:server:"+strconv.Itoa(int(server.ID)), int64(server.VolumeTotal))

	middleware.Success(ctx)
}
