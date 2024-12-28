package instanceController

import (
	"errors"
	"megrez/models"
	"megrez/services/database"
	"strings"
)

type patchReqStruct struct {
	CpuPatch    *cpuPatchStruct    `json:"cpuPatch"`
	GpuPatch    *gpuPatchStruct    `json:"gpuPatch"`
	MemoryPatch *MemoryPatchStruct `json:"memoryPatch"`
	VolumePatch *volumePatchStruct `json:"volumePatch"`
}

type cpuPatchStruct struct {
	CpuCount int `json:"cpuCount"`
}

type gpuPatchStruct struct {
	GpuCount int `json:"gpuCount"`
}

type MemoryPatchStruct struct {
	Memory string `json:"memory"`
}

type volumePatchStruct struct {
	OldBind bindStruct `json:"oldBind"`
	NewBind bindStruct `json:"newBind"`
}

func Patch(instance *models.Instances, gpuCount, volumeSize int, cpuOnly bool) (err error) {
	l.SetFunction("Patch")

	instance.Status = models.InstanceModifying
	result := database.DB.Save(&instance)
	if result.Error != nil {
		l.Error("save instance error: %v", result.Error)
		return result.Error
	}

	if gpuCount == instance.GpuCount && volumeSize == instance.VolumeSize && cpuOnly == instance.CpuOnly {
		instance.Status = models.InstanceStopped
		result = database.DB.Save(&instance)
		if result.Error != nil {
			l.Error("save instance error: %v", result.Error)
			return result.Error
		}
		return errors.New("no change")
	}

	server := models.Servers{
		ID: instance.ServerID,
	}
	result = database.DB.First(&server)
	if result.Error != nil {
		l.Error("query server error: %v", result.Error)
		return result.Error
	}

	oldVolumeName := instance.VolumeName
	if volumeSize != instance.VolumeSize {
		newVolumeName, err := patchVolume(server.IP, server.Port, server.Apikey, strings.Split(instance.VolumeName, "-")[0], volumeSize)
		if err != nil {
			l.Error("patch volume error: %v", err)
			return err
		}
		instance.VolumeName = newVolumeName
		instance.VolumeSize = volumeSize

		defer func(server models.Servers, instance *models.Instances, volumeName string) {
			err := deleteVolume(server.IP, server.Port, server.Apikey, volumeName, true)
			if err != nil {
				l.Error("delete volume error: %v", err)
			}
		}(server, instance, oldVolumeName)
	}

	if cpuOnly {
		err = patchCpuOnly(server.IP, server.Port, server.Apikey, instance.ContainerName, instance.VolumeName, oldVolumeName)
		if err != nil {
			l.Error("patch cpu only error: %v", err)
			return err
		}

		err = SetRootPassword(server.IP, server.Port, server.Apikey, instance.ContainerName, instance.SshPasswd)
		if err != nil {
			deleteInstance(server.IP, server.Port, server.Apikey, instance.ContainerName)
			if instance.VolumeName != "" {
				deleteVolume(server.IP, server.Port, server.Apikey, instance.VolumeName, false)
			}
			l.Error("set root password error: %v", err)
			return err
		}

		portBindings, err := GetPortForward(server.IP, server.Port, server.Apikey, instance.ContainerName)
		if err != nil {
			l.Error("get port forward error: %v", err)
			return err
		}

		instance.SshAddress = server.IP + ":" + portBindings["22"]
		instance.TensorBoardAddress = server.IP + ":" + portBindings["6007"]
		instance.JupyterAddress = server.IP + ":" + portBindings["8888"]
		instance.GrafanaAddress = server.IP + ":" + portBindings["3000"]

		instance.CpuOnly = true
		instance.GpuCount = 0
		instance.Status = models.InstanceRunning
		result = database.DB.Save(&instance)
		if result.Error != nil {
			l.Error("save instance error: %v", result.Error)
			return result.Error
		}
		return nil
	}

	err = patchGpu(server.IP, server.Port, server.Apikey,
		instance.ContainerName,
		server.CpuCountPerGpu, server.MemoryPerGpu,
		instance.VolumeName, oldVolumeName,
		gpuCount, instance.GpuCount,
	)
	if err != nil {
		l.Error("patch gpu error: %v", err)
		return err
	}

	err = SetRootPassword(server.IP, server.Port, server.Apikey, instance.ContainerName, instance.SshPasswd)
	if err != nil {
		deleteInstance(server.IP, server.Port, server.Apikey, instance.ContainerName)
		if instance.VolumeName != "" {
			deleteVolume(server.IP, server.Port, server.Apikey, instance.VolumeName, false)
		}
		l.Error("set root password error: %v", err)
		return err
	}

	portBindings, err := GetPortForward(server.IP, server.Port, server.Apikey, instance.ContainerName)
	if err != nil {
		l.Error("get port forward error: %v", err)
		return err
	}

	instance.SshAddress = server.IP + ":" + portBindings["22"]
	instance.TensorBoardAddress = server.IP + ":" + portBindings["6007"]
	instance.JupyterAddress = server.IP + ":" + portBindings["8888"]
	instance.GrafanaAddress = server.IP + ":" + portBindings["3000"]

	instance.CpuOnly = false
	instance.GpuCount = gpuCount
	instance.Status = models.InstanceRunning
	result = database.DB.Save(&instance)
	if result.Error != nil {
		l.Error("save instance error: %v", result.Error)
		return result.Error
	}

	return nil
}
