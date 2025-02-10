package models

import (
	"slices"
	"time"

	"gorm.io/gorm"
)

type Status int
type Action int

const (
	InstanceStatusFail    Status = -1
	InstanceStatusRunning Status = 0
	InstanceStatusPaused  Status = 1
	InstanceStatusStopped Status = 2

	InstanceStatusReady      Status = 3
	InstanceStatusStarting   Status = 4
	InstanceStatusStopping   Status = 5
	InstanceStatusPausing    Status = 6
	InstanceStatusRestarting Status = 7
	InstanceStatusModifying  Status = 8
	InstanceStatusDeleting   Status = 9
)

const (
	InstanceActionCreate  Action = 1
	InstanceActionStart   Action = 2
	InstanceActionPause   Action = 3
	InstanceActionStop    Action = 4
	InstanceActionRestart Action = 5
	InstanceActionModify  Action = 6
	InstanceActionDelete  Action = 7
)

var instanceIngStatus = []Status{InstanceStatusReady, InstanceStatusStarting, InstanceStatusStopping, InstanceStatusPausing, InstanceStatusRestarting, InstanceStatusModifying, InstanceStatusDeleting}

type Instances struct {
	ID uint `json:"id" gorm:"primary_key;autoIncrement;index"`

	UserID   uint `json:"user_id,omitempty" gorm:"not null"`
	ServerID uint `json:"server_id" gorm:"not null"`

	ImageName     string `json:"image_name" gorm:"type:varchar(255);not null"`
	ContainerName string `json:"container_name,omitempty" gorm:"type:varchar(255);not null"`
	CpuOnly       bool   `json:"cpu_only" gorm:"not null"`
	GpuCount      int    `json:"gpu_count" gorm:"not null"`
	VolumeName    string `json:"volume_name,omitempty" gorm:"type:varchar(255);not null"`
	VolumeSize    int    `json:"volume_size" gorm:"not null"`

	SshAddress         string `json:"ssh_address" gorm:"type:varchar(255)"`
	SshPasswd          string `json:"ssh_passwd" gorm:"type:varchar(255)"`
	JupyterAddress     string `json:"jupyter_address" gorm:"type:varchar(255)"`
	TensorBoardAddress string `json:"tensor_board_address" gorm:"type:varchar(255)"`
	GrafanaAddress     string `json:"grafana_address" gorm:"type:varchar(255)"`
	CodeServerAddress  string `json:"code_server_address" gorm:"type:varchar(255)"`
	Status             Status `json:"status" gorm:"not null"` // Detail in Constants
	FromAction         Action `json:"from_action"`

	Label string `json:"label" gorm:"type:varchar(255)"`

	CreatedAt time.Time      `json:"create_at"`
	UpdatedAt time.Time      `json:"update_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func InstanceIngStatusCheck(status Status) bool {
	if slices.Index(instanceIngStatus, status) != -1 {
		return true
	}
	return false
}
