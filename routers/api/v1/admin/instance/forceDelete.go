package instances

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"
	"megrez/services/dispatcher"
	"megrez/services/redis"
	"strconv"

	"github.com/kataras/iris/v12"
)

func forceDeleteHandler(ctx iris.Context) {
	l.SetFunction("forceDeleteHandler")

	id, err := ctx.Params().GetUint("id")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	instance := models.Instances{
		ID: id,
	}
	result := database.DB.First(&instance)
	if result.Error != nil {
		l.Error("get instance error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeInstanceDeleteError, iris.StatusInternalServerError)
		return
	}

	status := instance.Status
	if status != models.InstanceStatusFail {
		middleware.Error(ctx, middleware.CodeInstanceStatusError, iris.StatusBadRequest)
		return
	}

	if instance.FromAction == models.InstanceActionStop || instance.FromAction == models.InstanceActionPause || instance.FromAction == models.InstanceActionRestart {
		redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(instance.ServerID)), int64(instance.GpuCount))
	}
	if instance.FromAction != models.InstanceActionCreate {
		redis.RawDB.IncrBy(ctx, "remain_volume:server:"+strconv.Itoa(int(instance.ServerID)), int64(instance.VolumeSize+30))
	}

	dispatcherData := dispatcher.Data{
		Type:       dispatcher.Delete,
		Status:     status,
		InstanceID: instance.ID,
		Force:      true,
	}
	dispatcher.Push(instance.ServerID, dispatcherData)

	middleware.Success(ctx)

}
