package instanceController

import (
	"megrez/models"
	"megrez/services/database"
)

func Create(instance *models.Instances) (containerName, volumeName string, err error) {
	l.SetFunction("Create")

	server := models.Servers{
		ID: instance.ServerID,
	}
	result := database.DB.First(&server)
	if result.Error != nil {
		l.Error("query server error: %v", result.Error)
		return "", "", result.Error
	}

	if instance.VolumeSize != 0 {
		volumeName, err = createVolume(server.IP, server.Port, server.Apikey, instance.VolumeSize)
		if err != nil {
			l.Error("create volume error: %v", err)
			return "", "", err
		}
		instance.VolumeName = volumeName
	}

	containerName, err = createInstance(server.IP, server.Port, server.Apikey, instance.GpuCount, server.CpuCountPerGpu*instance.GpuCount, server.MemoryPerGpu*instance.GpuCount, instance.VolumeName, instance.ImageName)
	if err != nil {
		l.Error("create instance error: %v", err)
		deleteInstance(server.IP, server.Port, server.Apikey, containerName)
		if instance.VolumeName != "" {
			deleteVolume(server.IP, server.Port, server.Apikey, instance.VolumeName, false)
		}
		return "", "", err
	}
	instance.ContainerName = containerName

	err = SetRootPassword(server.IP, server.Port, server.Apikey, containerName, instance.SshPasswd)
	if err != nil {
		deleteInstance(server.IP, server.Port, server.Apikey, containerName)
		if instance.VolumeName != "" {
			deleteVolume(server.IP, server.Port, server.Apikey, instance.VolumeName, false)
		}
		l.Error("set root password error: %v", err)
		return "", "", err
	}

	go SetJupterPassword(server.IP, server.Port, server.Apikey, containerName, instance.SshPasswd)

	portBindings, err := GetPortForward(server.IP, server.Port, server.Apikey, containerName)
	if err != nil {
		deleteInstance(server.IP, server.Port, server.Apikey, containerName)
		if instance.VolumeName != "" {
			deleteVolume(server.IP, server.Port, server.Apikey, instance.VolumeName, false)
		}
		l.Error("get port forward error: %v", err)
		return "", "", err
	}

	instance.SshAddress = server.IP + ":" + portBindings["22"]
	instance.TensorBoardAddress = server.IP + ":" + portBindings["6007"]
	instance.JupyterAddress = server.IP + ":" + portBindings["8888"]
	instance.GrafanaAddress = server.IP + ":" + portBindings["3000"]

	instance.Status = 0
	result = database.DB.Save(&instance)

	return containerName, volumeName, nil
}
