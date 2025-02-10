package dispatcher

import (
	"context"
	"megrez/models"
	"megrez/services/database"
	"megrez/services/instanceController"
	"megrez/services/redis"
	"strconv"
)

func add(serverID uint, data Data) (err error) {
	lc := l.Clone()
	lc.SetFunction("add")

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

	containerName, volumeName, err := instanceController.Create(&instance)
	if err != nil {
		lc.Error("create instance error: %v", err)
		ctx := context.Background()
		redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(serverID)), int64(instance.GpuCount))
		redis.RawDB.IncrBy(ctx, "remain_volume:server:"+strconv.Itoa(int(serverID)), int64(instance.VolumeSize+30))
		database.DB.Model(&instance).Update("status", models.InstanceStatusFail).Update("from_action", models.InstanceActionCreate)
		return
	}

	server.GpuUsed += instance.GpuCount
	server.VolumeUsed += instance.VolumeSize + 30
	result = database.DB.Save(&server)
	if result.Error != nil {
		lc.Error("save server error: %v", result.Error)
		return result.Error
	}

	if volumeName != "" {
		lc.Info("create volume success: %v", volumeName)
	}
	lc.Info("create instance success: %v", containerName)

	return
}
