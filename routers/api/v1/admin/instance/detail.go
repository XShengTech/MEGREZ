package instances

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

	instance := models.Instances{
		ID: id,
	}
	result := database.DB.Select("id", "user_id", "server_id", "cpu_only", "gpu_count", "volume_size", "ssh_address", "ssh_passwd", "tensor_board_address", "grafana_address", "status", "image_name", "created_at").First(&instance)
	if result.Error != nil {
		l.Error("detail instance error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeInstanceDetailError, iris.StatusInternalServerError)
		return
	}

	server := models.Servers{
		ID: instance.ServerID,
	}
	result = database.DB.Select("name", "gpu_type", "gpu_num").First(&server)
	if result.Error != nil {
		l.Error("detail server error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeInstanceDetailError, iris.StatusInternalServerError)
		return
	}

	redis.RawDB.Get(ctx, "remain_gpu:server:"+strconv.Itoa(int(instance.ServerID))).Scan(&server.GpuUsed)
	server.GpuUsed = server.GpuNum - server.GpuUsed

	res := instanceStruct{
		Instances: instance,
		GpuType:   server.GpuType,
		GpuNum:    server.GpuNum,
		GpuUsed:   server.GpuUsed,
	}

	middleware.Result(ctx, res)
}
