package instances

import (
	"megrez/models"
	"megrez/routers/api/v1/middleware"
	"megrez/services/database"
	"megrez/services/dispatcher"
	"megrez/services/instanceController"
	"megrez/services/redis"
	"strconv"

	"github.com/kataras/iris/v12"
)

type controlStruct struct {
	Action instanceController.Action `json:"action"` // 1: start, 2: pause , 3: stop, 4: restart
}

func controlHandler(ctx iris.Context) {
	id, err := ctx.Params().GetUint("id")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	var req controlStruct
	err = ctx.ReadJSON(&req)
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	instance := models.Instances{
		ID: id,
	}
	result := database.DB.Where("id = ?", id).Where("user_id = ?", userId).First(&instance)
	if result.Error != nil {
		l.Error("detail instance error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeInstanceDetailError, iris.StatusInternalServerError)
		return
	}

	if models.InstanceIngStatusCheck(instance.Status) {
		middleware.Error(ctx, middleware.CodeInstanceStatusError, iris.StatusBadRequest)
		return
	}

	if req.Action == instanceController.ActionStop && instance.Status != models.InstanceStatusRunning && instance.Status != models.InstanceStatusPaused {
		middleware.Error(ctx, middleware.CodeInstanceStatusError, iris.StatusBadRequest)
		return
	}

	if req.Action == instanceController.ActionPause && instance.Status != models.InstanceStatusRunning {
		middleware.Error(ctx, middleware.CodeInstanceStatusError, iris.StatusBadRequest)
		return
	}

	if req.Action == instanceController.ActionStart && instance.Status != models.InstanceStatusStopped && instance.Status != models.InstanceStatusPaused {
		middleware.Error(ctx, middleware.CodeInstanceStatusError, iris.StatusBadRequest)
		return
	}

	server := models.Servers{
		ID: instance.ServerID,
	}
	result = database.DB.First(&server)
	if result.Error != nil {
		l.Error("query server error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeServerQueryError, iris.StatusInternalServerError)
		return
	}

	status := instance.Status
	if status == models.InstanceStatusStopped && (req.Action == instanceController.ActionStart || req.Action == instanceController.ActionRestart) {
		remainGpu, err := redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(server.ID)), int64(-instance.GpuCount)).Result()
		if err != nil {
			l.Error("incrby gpu num error: %v", err)
			middleware.Error(ctx, middleware.CodeServerQueryError, iris.StatusInternalServerError)
			return
		}

		if remainGpu < 0 {
			redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(server.ID)), int64(instance.GpuCount))
			middleware.Error(ctx, middleware.CodeResourceInsufficient, iris.StatusBadRequest)
			return
		}
	}

	switch req.Action {
	case instanceController.ActionStart:
		result = database.DB.Model(&instance).Update("status", models.InstanceStatusStarting)
		if result.Error != nil {
			l.Error("update instance status error: %v", result.Error)
			middleware.Error(ctx, middleware.CodeInstanceStartError, iris.StatusInternalServerError)
			return
		}

	case instanceController.ActionPause:
		result = database.DB.Model(&instance).Update("status", models.InstanceStatusPausing)
		if result.Error != nil {
			l.Error("update instance status error: %v", result.Error)
			middleware.Error(ctx, middleware.CodeInstancePauseError, iris.StatusInternalServerError)
			return
		}

	case instanceController.ActionStop:
		result = database.DB.Model(&instance).Update("status", models.InstanceStatusStopping)
		if result.Error != nil {
			l.Error("update instance status error: %v", result.Error)
			middleware.Error(ctx, middleware.CodeInstanceStopError, iris.StatusInternalServerError)
			return
		}
		redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(server.ID)), int64(instance.GpuCount))

	case instanceController.ActionRestart:
		result = database.DB.Model(&instance).Update("status", models.InstanceStatusRestarting)
		if result.Error != nil {
			l.Error("update instance status error: %v", result.Error)
			middleware.Error(ctx, middleware.CodeInstanceRestartError, iris.StatusInternalServerError)
			return
		}

	default:
	}

	dispatcherData := dispatcher.Data{
		Type:       dispatcher.Control,
		InstanceID: instance.ID,
		Status:     status,
		Action:     req.Action,
	}
	dispatcher.Push(instance.ServerID, dispatcherData)

	middleware.Success(ctx)
}
