package system

import (
	"megrez/libs/crypto"
	"megrez/models"
	"megrez/services/config"
	"megrez/services/database"
)

const imagesKey = "images"

func systemInit() (err error) {
	l.SetFunction("systemInit")

	salt := crypto.Hex(32)
	config.SetSystemSalt(salt)
	err = config.Save()
	if err != nil {
		l.Error("Save config failed, Error: %v", err)
		return
	}

	user := models.Users{
		Username: "admin",
		Email:    "admin@gpuManager.com",
		Role:     3,
		Verify:   true,
	}
	result := database.DB.Create(&user)
	if result.Error != nil {
		l.Error("Create admin user failed, Error: %v", result.Error)
		return
	}
	user.Password = user.PasswordHash("admin123456")
	result = database.DB.Save(&user)
	if result.Error != nil {
		l.Error("Save admin user failed, Error: %v", result.Error)
		return
	}

	st := models.System{
		Key:   imagesKey,
		Value: "{}",
	}
	result = database.DB.Create(&st)
	if result.Error != nil {
		l.Error("Create system failed, Error: %v", result.Error)
		return
	}

	l.Info("System init success")

	return
}
