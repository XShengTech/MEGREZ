package instanceController

import (
	"encoding/json"
	"errors"
	"megrez/libs/request"
	"megrez/models"
	"megrez/services/database"
	"strconv"
)

func Stop(instance *models.Instances) (err error) {
	l.SetFunction("Stop")

	instance.Status = models.InstanceStopping
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

	err = stopInstance(server.IP, server.Port, server.Apikey, instance.ContainerName)
	if err != nil {
		l.Error("stop instance error: %v", err)
		return err
	}

	instance.Status = models.InstanceStopped
	result = database.DB.Save(&instance)
	if result.Error != nil {
		l.Error("save instance error: %v", result.Error)
		return result.Error
	}

	return nil
}

func stopInstance(ip string, port int, apikey, containerName string) (err error) {
	l.SetFunction("stopInstance")

	c := request.NewRequest().Patch().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + instancePrefix + "/" + containerName + instanceStop).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez")
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("stop instance error: %v", c.GetBody())
		return errors.New("stop instance request error")
	}

	var res resStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return err
	}

	if res.Code != 200 {
		l.Error("stop instance code: %d, error: %v", res.Code, res.Msg)
		return errors.New(res.Msg)
	}

	return nil
}
