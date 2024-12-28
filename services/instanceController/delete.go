package instanceController

import (
	"megrez/models"
	"megrez/services/database"
)

func Delete(instance *models.Instances) (err error) {
	l.SetFunction("Delete")

	instance.Status = models.InstanceDeleting
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

	err = deleteInstance(server.IP, server.Port, server.Apikey, instance.ContainerName)
	if err != nil {
		l.Error("delete instance error: %v", err)
		return err
	}

	if instance.VolumeName != "" {
		err = deleteVolume(server.IP, server.Port, server.Apikey, instance.VolumeName, false)
		if err != nil {
			l.Error("delete volume error: %v", err)
		}
	}

	result = database.DB.Delete(&instance)
	if result.Error != nil {
		l.Error("delete instance error: %v", result.Error)
		return result.Error
	}

	return nil
}
