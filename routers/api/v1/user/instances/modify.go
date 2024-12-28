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

type modifyReqStruct struct {
	CpuOnly    bool `json:"cpu_only"`
	GpuCount   *int `json:"gpu_count"`
	VolumeSize *int `json:"volume_size"`
}

func modifyHandler(ctx iris.Context) {
	l.SetFunction("modifyHandler")

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

	var req modifyReqStruct
	err = ctx.ReadJSON(&req)
	if err != nil {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
		return
	}

	if req.GpuCount != nil {
		if *req.GpuCount < 0 {
			middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
			return
		}
	}

	if req.VolumeSize != nil {
		if *req.VolumeSize < 50 {
			middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
			return
		}
	}

	instance := models.Instances{
		ID: id,
	}
	result := database.DB.Where(&instance).Where("user_id = ?", userId).First(&instance)
	if result.Error != nil {
		l.Error("query instance error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeInstanceQueryError, iris.StatusInternalServerError)
		return
	}

	if instance.Status != models.InstanceStopped {
		middleware.Error(ctx, middleware.CodeInstanceStatusError, iris.StatusBadRequest)
		return
	}

	if req.CpuOnly == instance.CpuOnly && req.CpuOnly {
		middleware.Error(ctx, middleware.CodeBadRequest, iris.StatusBadRequest)
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

	if req.GpuCount != nil {
		remainGpu, err := redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(server.ID)), int64(-*req.GpuCount)).Result()
		if err != nil {
			l.Error("incrby gpu num error: %v", err)
			middleware.Error(ctx, middleware.CodeServerQueryError, iris.StatusInternalServerError)
			return
		}

		if remainGpu < 0 {
			redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(server.ID)), int64(*req.GpuCount))
			middleware.Error(ctx, middleware.CodeResourceInsufficient, iris.StatusBadRequest)
			return
		}
	}
	if req.VolumeSize != nil {
		remainVolume, err := redis.RawDB.IncrBy(ctx, "remain_volume:server:"+strconv.Itoa(int(server.ID)), int64(instance.VolumeSize-*req.VolumeSize)).Result()
		if err != nil {
			l.Error("incrby volume size error: %v", err)
			middleware.Error(ctx, middleware.CodeServerQueryError, iris.StatusInternalServerError)
			return
		}
		if remainVolume < 0 {
			redis.RawDB.IncrBy(ctx, "remain_volume:server:"+strconv.Itoa(int(server.ID)), int64(*req.VolumeSize-instance.VolumeSize))
			middleware.Error(ctx, middleware.CodeResourceInsufficient, iris.StatusBadRequest)
			return
		}
	}

	status := instance.Status
	result = database.DB.Model(&instance).Update("status", models.InstanceModifying)
	if result.Error != nil {
		l.Error("update instance status error: %v", result.Error)
		middleware.Error(ctx, middleware.CodeServerSaveError, iris.StatusInternalServerError)
		return
	}

	dispatcherData := dispatcher.Data{
		Type:       dispatcher.Modify,
		InstanceID: instance.ID,
		Status:     status,

		CpuOnly:    req.CpuOnly,
		GpuCount:   req.GpuCount,
		VolumeSize: req.VolumeSize,
	}
	dispatcher.Push(instance.ServerID, dispatcherData)

	middleware.Success(ctx)
}
