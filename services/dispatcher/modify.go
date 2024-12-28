package dispatcher

import (
	"context"
	"errors"
	"megrez/models"
	"megrez/services/database"
	"megrez/services/instanceController"
	"megrez/services/redis"
	"strconv"
)

type modifyDataStruct struct {
	CpuOnly    bool `json:"cpu_only"`
	GpuCount   *int `json:"gpu_count"`
	VolumeSize *int `json:"volume_size"`
}

func modify(serverID uint, data Data) (err error) {
	lc := l.Clone()
	lc.SetFunction("modify")

	server := models.Servers{
		ID: serverID,
	}
	result := database.DB.First(&server)
	if result.Error != nil {
		lc.Error("query server error: %v", result.Error)
		return
	}

	instance := models.Instances{
		ID: data.InstanceID,
	}
	result = database.DB.First(&instance)
	if result.Error != nil {
		lc.Error("query instance error: %v", result.Error)
		return result.Error
	}

	if data.Status != models.InstanceStopped {
		lc.Error("instance status error")
		return errors.New("instance status error")
	}

	if data.CpuOnly == instance.CpuOnly && data.CpuOnly {
		lc.Error("instance already cpu_only mode")
		return errors.New("instance already cpu_only mode")
	}

	oldVolumeSize := instance.VolumeSize
	gpuCount := instance.GpuCount
	volumeSize := instance.VolumeSize
	if data.GpuCount != nil {
		gpuCount = *data.GpuCount
	}
	if data.VolumeSize != nil {
		volumeSize = *data.VolumeSize
	}

	err = instanceController.Patch(&instance, gpuCount, volumeSize, data.CpuOnly)
	if err != nil {
		lc.Error("patch instance error: %v", err)
		ctx := context.Background()
		redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(serverID)), int64(gpuCount))
		redis.RawDB.IncrBy(ctx, "remain_volume:server:"+strconv.Itoa(int(serverID)), int64(volumeSize-oldVolumeSize))
		return
	}

	if data.Status == models.InstanceStopped {
		server.GpuUsed += instance.GpuCount
	}
	server.VolumeUsed += instance.VolumeSize - oldVolumeSize

	result = database.DB.Save(&server)
	if result.Error != nil {
		lc.Error("save server error: %v", result.Error)
		return result.Error
	}

	lc.Info("modify instance success: %d", instance.ID)

	return
}
