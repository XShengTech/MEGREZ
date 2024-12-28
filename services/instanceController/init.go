package instanceController

import (
	"megrez/services/logger"

	_logger "megrez/libs/logger"
)

var l *_logger.LoggerStruct

func InitInstanceController() {
	l = logger.Logger.Clone()
	l.SetModel("instanceController")
}

type bindStruct struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
}

type resStruct struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data map[string]any `json:"data"`
}
