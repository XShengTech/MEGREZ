package dispatcher

import (
	"context"
	"megrez/models"
	"megrez/services/database"
	"megrez/services/logger"
	"megrez/services/redis"
	"strconv"
)

func Init() {
	l = logger.Logger.Clone()
	l.SetModel("dispatcher")
	l.SetFunction("Init")

	status = make(map[int]bool)
	// run()

	var servers []models.Servers
	result := database.DB.Select("id", "gpu_num", "gpu_used", "volume_total", "volume_used").Find(&servers)
	if result.Error != nil {
		l.Error("query server error: %v", result.Error)
		return
	}

	ctx := context.Background()

	for _, server := range servers {
		status[int(server.ID)] = false
		go run(server.ID, false)

		_, err := redis.RawDB.Get(ctx, "remain_gpu:server:"+strconv.Itoa(int(server.ID))).Result()
		if err != nil {
			l.Info("init server %d remain gpu: %d", server.ID, server.GpuNum-server.GpuUsed)
			redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(server.ID)), int64(server.GpuNum-server.GpuUsed))
		}

		_, err = redis.RawDB.Get(ctx, "remain_volume:server:"+strconv.Itoa(int(server.ID))).Result()
		if err != nil {
			l.Info("init server %d remain volume: %d", server.ID, server.VolumeTotal-server.VolumeUsed)
			redis.RawDB.IncrBy(ctx, "remain_volume:server:"+strconv.Itoa(int(server.ID)), int64(server.VolumeTotal-server.VolumeUsed))
		}
	}
}
