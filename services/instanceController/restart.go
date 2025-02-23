package instanceController

import (
	"encoding/json"
	"errors"
	"megrez/libs/request"
	"megrez/models"
	"megrez/services/database"
	"strconv"
)

func Restart(instance *models.Instances) (err error) {
	l.SetFunction("Restart")

	if instance.Status == models.InstanceStatusStopped || instance.Status == models.InstanceStatusPaused {
		instance.Status = models.InstanceStatusStarting
	} else {
		instance.Status = models.InstanceStatusRestarting
	}
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

	err = restartInstance(server.IP, server.Port, server.Apikey, instance.ContainerName)
	if err != nil {
		l.Error("restart instance error: %v", err)
		return err
	}

	err = SetRootPassword(server.IP, server.Port, server.Apikey, instance.ContainerName, instance.SshPasswd)
	if err != nil {
		l.Error("set root password error: %v", err)
		return err
	}

	go func() {
		SetJupyterPassword(server.IP, server.Port, server.Apikey, instance.ContainerName, instance.SshPasswd)
		SetCodeServerPassword(server.IP, server.Port, server.Apikey, instance.ContainerName, instance.SshPasswd)
	}()

	portBindings, err := GetPortForward(server.IP, server.Port, server.Apikey, instance.ContainerName)
	if err != nil {
		l.Error("get port forward error: %v", err)
		return err
	}

	instance.SshAddress = server.IP + ":" + portBindings["22"]
	instance.TensorBoardAddress = server.IP + ":" + portBindings["6007"]
	instance.JupyterAddress = server.IP + ":" + portBindings["80"] + "/jupyter"
	instance.GrafanaAddress = server.IP + ":" + portBindings["80"] + "/monitor/public-dashboards/2c510f203876465ba76617510ce3e219"
	instance.CodeServerAddress = server.IP + ":" + portBindings["80"] + "/code-server/"

	instance.Status = models.InstanceStatusRunning
	result = database.DB.Save(&instance)
	if result.Error != nil {
		l.Error("save instance error: %v", result.Error)
		return result.Error
	}

	return nil
}

func restartInstance(ip string, port int, apikey, containerName string) (err error) {
	l.SetFunction("restartInstance")

	c := request.NewRequest().Patch().
		SetUrl("http://" + ip + ":" + strconv.Itoa(port) + apiPrefix + instancePrefix + "/" + containerName + instanceRestart).
		SetAuthorization("Bearer " + apikey).
		SetUserAgent("megrez")
	c.Do()

	if c.GetStatusCode() != 200 {
		l.Error("restart instance error: %v", c.GetBody())
		return errors.New("restart instance request error")
	}

	var res resStruct
	err = json.Unmarshal(c.GetBody(), &res)
	if err != nil {
		l.Error("unmarshal response data error: %v", err)
		return err
	}

	if res.Code != 200 {
		l.Error("restart instance error: %v", res.Msg)
		return errors.New(res.Msg)
	}

	return nil
}
