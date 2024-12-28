package servers

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"
	"megrez/services/redis"
	"strconv"

	"github.com/kataras/iris/v12"
)

func detailHandler(ctx iris.Context) {
	l.SetFunction("detailHandler")

	id, err := ctx.Params().GetUint("id")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	server := models.Servers{
		ID: id,
	}
	result := database.DB.First(&server)
	if result.Error != nil {
		l.Error("detail server error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeAdminServerDetailError, iris.StatusInternalServerError)
		return
	}

	redis.RawDB.Get(ctx, "remain_gpu:server:"+strconv.Itoa(int(id))).Scan(&server.GpuUsed)
	server.GpuUsed = server.GpuNum - server.GpuUsed

	redis.RawDB.Get(ctx, "remain_volume:server:"+strconv.Itoa(int(id))).Scan(&server.VolumeUsed)
	server.VolumeUsed = server.VolumeTotal - server.VolumeUsed

	middleware.Result(ctx, server)
}
