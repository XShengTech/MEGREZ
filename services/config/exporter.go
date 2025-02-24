package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

func InitConfig(path string) {
	l.SetFunction("InitConfig")

	configPath = path
	c, err := os.ReadFile(path)
	if err != nil {
		l.Fatal("Failed to read config file", err)
	}
	err = yaml.Unmarshal(c, &config)
	if err != nil {
		l.Fatal("Failed to unmarshal config file", err)
	}
	l.SetLevel(config.GetLogLevel())

	if config.System.BaseUrl == "" {
		config.System.BaseUrl = "http://" + config.Http.Host + ":" + strconv.Itoa(config.Http.Port)
	}
}

func GetDatabase() databaseStruct {
	return config.GetDatabase()
}

func GetHttpAddress() string {
	return config.GetHttpAddress()
}

func GetRedis() redisStruct {
	return config.GetRedis()
}

func GetSmtp() smtpStruct {
	return config.GetSmtp()
}

func GetLogLevel() string {
	return config.GetLogLevel()
}

func GetLogFile() string {
	return config.GetLogFile()
}

func GetSystemBaseUrl() string {
	return config.GetSystemBaseUrl()
}

func GetSystemSalt() string {
	return config.GetSystemSalt()
}

func SetSystemSalt(salt string) {
	config.SetSystemSalt(salt)
}

func GetSystemVerify() bool {
	return config.GetSystemVerify()
}

func GetSystemMountDir() string {
	return config.GetSystemMountDir()
}

func Save() error {
	return config.Save()
}
