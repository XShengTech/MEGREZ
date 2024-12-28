package models

import (
	"slices"
	"time"

	"gorm.io/gorm"
)

type Status int

const (
	InstanceFail    Status = -1
	InstanceRunning Status = 0
	InstancePaused  Status = 1
	InstanceStopped Status = 2

	InstanceReady      Status = 3
	InstanceStarting   Status = 4
	InstanceStopping   Status = 5
	InstancePausing    Status = 6
	InstanceRestarting Status = 7
	InstanceModifying  Status = 8
	InstanceDeleting   Status = 9
)

var instanceIngStatus = []Status{InstanceReady, InstanceStarting, InstanceStopping, InstancePausing, InstanceRestarting, InstanceModifying, InstanceDeleting}

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

	SshAddress         string `json:"ssh_address" gorm:"type:varchar(255);not null"`
	SshPasswd          string `json:"ssh_passwd" gorm:"type:varchar(255);not null"`
	JupyterAddress     string `json:"jupyter_address" gorm:"type:varchar(255);not null"`
	TensorBoardAddress string `json:"tensor_board_address" gorm:"type:varchar(255);not null"`
	GrafanaAddress     string `json:"grafana_address" gorm:"type:varchar(255);not null"`
	Status             Status `json:"status" gorm:"not null"` // -1: fail, 0: running, 1: paused, 2: stopped, 3: readying

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
