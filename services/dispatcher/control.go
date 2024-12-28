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

func control(serverID uint, data Data) (err error) {
	lc := l.Clone()
	lc.SetFunction("control")

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

	if models.InstanceIngStatusCheck(data.Status) {
		lc.Error("instance status error")
		return errors.New("instance status error")
	}

	if data.Action == instanceController.ActionPause && data.Status != models.InstanceRunning {
		lc.Error("instance status error")
		return errors.New("instance status error")
	}

	if data.Action == instanceController.ActionStart && data.Status != models.InstanceStopped && data.Status != models.InstancePaused {
		lc.Error("instance status error")
		return errors.New("instance status error")
	}

	switch data.Action {
	case instanceController.ActionStart:
		if data.Status == models.InstancePaused {
			err = instanceController.Continue(&instance)
			if err != nil {
				lc.Error("instance continue error: %v", err)
				return
			}
		} else if data.Status == models.InstanceStopped {
			err = instanceController.Restart(&instance)
			if err != nil {
				lc.Error("instance restart error: %v", err)
				ctx := context.Background()
				redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(serverID)), int64(instance.GpuCount))
				return
			}
		}

		if data.Status == models.InstanceStopped && instance.Status == models.InstanceRunning {
			server.GpuUsed += instance.GpuCount
			result = database.DB.Save(&server)
			if result.Error != nil {
				lc.Error("save server error: %v", result.Error)
				return result.Error
			}
		}

		lc.Info("start instance success: %v", instance.ID)

	case instanceController.ActionPause:
		err = instanceController.Pause(&instance)
		if err != nil {
			lc.Error("instance pause error: %v", err)
			return
		}

		lc.Info("pause instance success: %v", instance.ID)

	case instanceController.ActionStop:
		err = instanceController.Stop(&instance)
		if err != nil {
			lc.Error("instance stop error: %v", err)
			return
		}

		if (data.Status == models.InstanceRunning || data.Status == models.InstancePaused) && instance.Status == models.InstanceStopped {
			server.GpuUsed -= instance.GpuCount
			result = database.DB.Save(&server)
			if result.Error != nil {
				lc.Error("save server error: %v", result.Error)
				return result.Error
			}
		}

		lc.Info("stop instance success: %v", instance.ID)

	case instanceController.ActionRestart:
		err = instanceController.Restart(&instance)
		if err != nil {
			lc.Error("instance restart error: %v", err)
			return
		}

		if data.Status == models.InstanceStopped && instance.Status == models.InstanceRunning {
			server.GpuUsed += instance.GpuCount
			result = database.DB.Save(&server)
			if result.Error != nil {
				lc.Error("save server error: %v", result.Error)
				return result.Error
			}
		}

		lc.Info("restart instance success: %v", instance.ID)

	default:
	}

	return
}
