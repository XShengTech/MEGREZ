package system

import (
	_logger "megrez/libs/logger"
	"megrez/services/config"
	"megrez/services/logger"
)

var l *_logger.LoggerStruct

func Check() (err error) {
	l = logger.Logger.Clone()
	l.SetModel("system")
	l.SetFunction("Check")

	salt := config.GetSystemSalt()
	if salt == "" {
		l.Info("System not init, initializing")
		systemInit()
	} else {
		l.Debug("System already init")
	}

	return nil
}
