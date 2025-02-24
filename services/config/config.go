package config

import "megrez/services/logger"

var l = logger.Logger.Clone()
var configPath string

type configStruct struct {
	Http     httpStruct     `yaml:"http,omitempty"`
	Database databaseStruct `yaml:"database,omitempty"`
	Redis    redisStruct    `yaml:"redis,omitempty"`
	Smtp     smtpStruct     `yaml:"smtp,omitempty"`
	Log      logStruct      `yaml:"log,omitempty"`
	System   systemStruct   `yaml:"system,omitempty"`
}

type httpStruct struct {
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
}

type databaseStruct struct {
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Database string `yaml:"database,omitempty"`
	SSLMode  bool   `yaml:"ssl,omitempty"`
}

type redisStruct struct {
	Host             string `yaml:"host,omitempty"`
	Port             int    `yaml:"port,omitempty"`
	Password         string `yaml:"password,omitempty"`
	Database         int    `yaml:"database,omitempty"`
	SentinelPassword string `yaml:"sentinel_password,omitempty"`
}

type smtpStruct struct {
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Password string `yaml:"password,omitempty"`
	User     string `yaml:"user,omitempty"`
	SSL      bool   `yaml:"ssl,omitempty"`
}

type logStruct struct {
	Level string `yaml:"level,omitempty"`
	File  string `yaml:"file,omitempty"`
}

type systemStruct struct {
	BaseUrl  string `yaml:"base_url,omitempty"`
	Salt     string `yaml:"salt,omitempty"`
	Verify   bool   `yaml:"verify,omitempty"`
	MountDir string `yaml:"mount_dir,omitempty"`
}

var config = configStruct{
	Http: httpStruct{
		Host: "0.0.0.0",
		Port: 8080,
	},
	Database: databaseStruct{
		Host:     "localhost",
		Port:     5432,
		Username: "GpuManager",
		Password: "GpuManager",
		Database: "GpuManager",
	},
	Redis: redisStruct{
		Host:     "localhost",
		Port:     6379,
		Password: "GpuManager",
		Database: 0,
	},
	Smtp: smtpStruct{},
	Log: logStruct{
		Level: "DEBUG",
		File:  "data/logs/backend.log",
	},
	System: systemStruct{
		Salt:   "",
		Verify: false,
	},
}
