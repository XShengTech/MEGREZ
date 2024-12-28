package dispatcher

import (
	"context"
	"megrez/models"
	"megrez/services/database"
	"megrez/services/instanceController"
	"megrez/services/redis"
	"strconv"
)

func delete(serverID uint, data Data) (err error) {
	lc := l.Clone()
	lc.SetFunction("delete")

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

	err = instanceController.Delete(&instance)
	if err != nil {
		lc.Error("delete instance error: %v", err)
		ctx := context.Background()
		if data.Status == models.InstanceRunning || data.Status == models.InstancePaused {
			redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(serverID)), int64(-instance.GpuCount))
		}
		redis.RawDB.IncrBy(ctx, "remain_volume:server:"+strconv.Itoa(int(serverID)), int64(-instance.VolumeSize-30))
		return
	}

	server.VolumeUsed -= instance.VolumeSize + 30
	if data.Status == models.InstanceRunning || data.Status == models.InstancePaused {
		server.GpuUsed -= instance.GpuCount
	}
	result = database.DB.Save(&server)
	if result.Error != nil {
		lc.Error("save server error: %v", result.Error)
		return result.Error
	}

	lc.Info("delete instance success: %v", instance.ID)

	return
}
