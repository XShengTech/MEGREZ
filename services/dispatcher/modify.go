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

	if data.Status != models.InstanceStatusStopped {
		lc.Error("instance status error")
		return errors.New("instance status error")
	}

	hasChanges := false
	if data.CpuOnly != instance.CpuOnly {
		hasChanges = true
	}
	if data.GpuCount != nil && *data.GpuCount != instance.GpuCount {
		hasChanges = true
	}
	if data.VolumeSize != nil && *data.VolumeSize != instance.VolumeSize {
		hasChanges = true
	}
	if !hasChanges {
		lc.Error("no changes")
		return errors.New("no changes")
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
		ctx := context.Background()
		redis.RawDB.IncrBy(ctx, "remain_volume:server:"+strconv.Itoa(int(serverID)), int64(volumeSize-oldVolumeSize))
		database.DB.Model(&instance).Update("status", models.InstanceStatusFail).Update("from_action", models.InstanceActionModify)
		lc.Error("patch instance error: %v", err)
		return
	}

	if data.Status == models.InstanceStatusStopped {
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
