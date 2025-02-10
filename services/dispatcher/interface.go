package dispatcher

import (
	_logger "megrez/libs/logger"
	"megrez/models"
	"megrez/services/instanceController"
	"sync"
)

var l *_logger.LoggerStruct
var status map[int]bool
var mx sync.Mutex

type Data struct {
	Type       Type          `json:"type"` // 1:Add 2:Contorl 3:Modify 4:Delete
	InstanceID uint          `json:"instance_id"`
	Status     models.Status `json:"status"`

	CpuOnly    bool `json:"cpu_only,omitempty"`
	GpuCount   *int `json:"gpu_count,omitempty"`
	VolumeSize *int `json:"volume_size,omitempty"`

	Action instanceController.Action `json:"action,omitempty"`
	Force  bool                      `json:"force,omitempty"`
}

type Type int

const (
	Add Type = iota + 1
	Control
	Modify
	Delete
)
