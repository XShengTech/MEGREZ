package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_logger "megrez/libs/logger"
	"megrez/services/config"
	"megrez/services/logger"

	__logger "gorm.io/gorm/logger"
)

var l *_logger.LoggerStruct

var DB *gorm.DB

func Connect() (err error) {
	l = logger.Logger.Clone()
	l.SetModel("database")
	l.SetFunction("Connect")

	c := config.GetDatabase()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", c.Host, c.Username, c.Password, c.Database, c.Port)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      __logger.Default.LogMode(__logger.Silent),
	})
	if err != nil {
		l.Fatal("Failed to connect to database", err)
	} else {
		l.Info("Database Connect Success")
	}

	autoMigrate()
	return
}
