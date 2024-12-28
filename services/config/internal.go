package config

import (
	"errors"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

func (c *configStruct) GetDatabase() databaseStruct {
	return c.Database
}

func (c *configStruct) GetHttpAddress() string {
	return c.Http.Host + ":" + strconv.Itoa(c.Http.Port)
}

func (c *configStruct) GetRedis() redisStruct {
	return c.Redis
}

func (c *configStruct) GetLogLevel() string {
	return c.Log.Level
}

func (c *configStruct) GetLogFile() string {
	return c.Log.File
}

func (c *configStruct) GetSystemSalt() string {
	return c.System.Salt
}

func (c *configStruct) SetSystemSalt(salt string) {
	c.System.Salt = salt
}

func (c *configStruct) GetSystemVerify() bool {
	return c.System.Verify
}

func (c *configStruct) Save() error {
	return c.save()
}

func (c *configStruct) save() error {
	// Save config to file
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return errors.New("yaml marshal Error: " + err.Error())
	}

	f, err := os.OpenFile(configPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return errors.New("open file Error: " + err.Error())
	}

	_, err = f.Write(bytes)
	if err != nil {
		return errors.New("write file Error: " + err.Error())
	}

	return nil
}
