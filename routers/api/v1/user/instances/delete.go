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

func deleteHandler(ctx iris.Context) {
	l.SetFunction("deleteHandler")

	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	id, err := ctx.Params().GetUint("id")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	var instance models.Instances
	result := database.DB.Where("id = ?", id).Where("user_id = ?", userId).First(&instance)
	if result.Error != nil {
		l.Error("get instance error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeInstanceDeleteError, iris.StatusInternalServerError)
		return
	}

	status := instance.Status
	if models.InstanceIngStatusCheck(status) {
		middleware.Error(ctx, middleware.CodeInstanceStatusError, iris.StatusBadRequest)
		return
	}

	result = database.DB.Model(&instance).Update("status", models.InstanceDeleting)
	if result.Error != nil {
		l.Error("update instance status error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeServerSaveError, iris.StatusInternalServerError)
		return
	}

	if status == models.InstanceRunning || status == models.InstancePaused {
		redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(instance.ServerID)), int64(instance.GpuCount))
	}
	redis.RawDB.IncrBy(ctx, "remain_volume:server:"+strconv.Itoa(int(instance.ServerID)), int64(instance.VolumeSize+30))

	// TODO: Price calculation

	dispatcherData := dispatcher.Data{
		Type:       dispatcher.Delete,
		Status:     status,
		InstanceID: instance.ID,
	}
	dispatcher.Push(instance.ServerID, dispatcherData)

	middleware.Success(ctx)
}
