package instanceController

import (
	"encoding/json"
	"errors"
	"megrez/libs/request"
	"megrez/models"
	"megrez/services/database"
	"strconv"
)

func Pause(instance *models.Instances) (err error) {
	l.SetFunction("Pause")

	instance.Status = models.InstancePausing
	result := database.DB.Save(&instance)
	if result.Error != nil {
		l.Error("save instance error: %v", result.Error)
		return result.Error
	}

	server := models.Servers{
		ID: instance.ServerID,
	}
	result = database.DB.First(&server)
	if result.Error != nil {
		l.Error("query server error: %v", result.Error)
		return result.Error
	}

	err = pauseInstance(server.IP, server.Port, server.Apikey, instance.ContainerName)
	if err != nil {
		l.Error("pause instance error: %v", err)
		return err
	}

	instance.Status = models.InstancePaused
	result = database.DB.Save(&instance)
	if result.Error != nil {
		l.Error("save instance error: %v", result.Error)
		return result.Error
	}

	return nil
}

func pauseInstance(ip string, port int, apikey, containerName string) (err error) {
	l.SetFunction("pauseInstance")

	c := request.NewRequest().Patch().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + instancePrefix + "/" + containerName + instancePause).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez")
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("pause instance error: %v", c.GetBody())
		return errors.New("pause instance request error")
	}

	var res resStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return err
	}

	if res.Code != 200 {
		l.Error("pause instance error: %v", res.Msg)
		return errors.New(res.Msg)
	}

	return nil
}
