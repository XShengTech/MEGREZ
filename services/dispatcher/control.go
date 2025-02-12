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

	if data.Action == instanceController.ActionPause && data.Status != models.InstanceStatusRunning {
		lc.Error("instance status error")
		return errors.New("instance status error")
	}

	if data.Action == instanceController.ActionStart && data.Status != models.InstanceStatusStopped && data.Status != models.InstanceStatusPaused {
		lc.Error("instance status error")
		return errors.New("instance status error")
	}

	switch data.Action {
	case instanceController.ActionStart:
		if data.Status == models.InstanceStatusPaused {
			err = instanceController.Continue(&instance)
			if err != nil {
				database.DB.Model(&instance).Update("status", models.InstanceStatusFail).Update("from_action", models.InstanceActionStart)
				lc.Error("instance continue error: %v", err)
				return
			}
		} else if data.Status == models.InstanceStatusStopped {
			err = instanceController.Restart(&instance)
			if err != nil {
				ctx := context.Background()
				redis.RawDB.IncrBy(ctx, "remain_gpu:server:"+strconv.Itoa(int(serverID)), int64(instance.GpuCount))
				database.DB.Model(&instance).Update("status", models.InstanceStatusFail).Update("from_action", models.InstanceActionStart)
				lc.Error("instance restart error: %v", err)
				return
			}
		}

		if data.Status == models.InstanceStatusStopped && instance.Status == models.InstanceStatusRunning {
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
			database.DB.Model(&instance).Update("status", models.InstanceStatusFail).Update("from_action", models.InstanceActionPause)
			lc.Error("instance pause error: %v", err)
			return
		}

		lc.Info("pause instance success: %v", instance.ID)

	case instanceController.ActionStop:
		err = instanceController.Stop(&instance)
		if err != nil {
			database.DB.Model(&instance).Update("status", models.InstanceStatusFail).Update("from_action", models.InstanceActionStop)
			lc.Error("instance stop error: %v", err)
			return
		}

		if (data.Status == models.InstanceStatusRunning || data.Status == models.InstanceStatusPaused) && instance.Status == models.InstanceStatusStopped {
			server.GpuUsed -= instance.GpuCount
			result = database.DB.Save(&server)
			if result.Error != nil {
				database.DB.Model(&instance).Update("status", models.InstanceStatusFail)
				lc.Error("save server error: %v", result.Error)
				return result.Error
			}
		}

		lc.Info("stop instance success: %v", instance.ID)

	case instanceController.ActionRestart:
		err = instanceController.Restart(&instance)
		if err != nil {
			database.DB.Model(&instance).Update("status", models.InstanceStatusFail).Update("from_action", models.InstanceActionRestart)
			lc.Error("instance restart error: %v", err)
			return
		}

		if data.Status == models.InstanceStatusStopped && instance.Status == models.InstanceStatusRunning {
			server.GpuUsed += instance.GpuCount
			result = database.DB.Save(&server)
			if result.Error != nil {
				database.DB.Model(&instance).Update("status", models.InstanceStatusFail)
				lc.Error("save server error: %v", result.Error)
				return result.Error
			}
		}

		lc.Info("restart instance success: %v", instance.ID)

	default:
	}

	return
}
