package models

import (
	"time"

	"gorm.io/gorm"
)

type Servers struct {
	ID uint `json:"id" gorm:"primary_key;autoIncrement;index"`

	Name             string `json:"name" gorm:"type:varchar(255);not null"`
	IP               string `json:"ip,omitempty" gorm:"type:varchar(255);not null"`
	Port             int    `json:"port,omitempty" gorm:"not null"`
	Apikey           string `json:"apikey,omitempty" gorm:"type:varchar(255);not null"`
	GpuType          string `json:"gpu_type" gorm:"type:varchar(255);not null"`
	GpuNum           int    `json:"gpu_num" gorm:"not null"`
	GpuDriverVersion string `json:"gpu_driver_version" gorm:"type:varchar(255);not null"`
	GpuCudaVersion   string `json:"gpu_cuda_version" gorm:"type:varchar(255);not null"`

	CpuCountPerGpu int     `json:"cpu_count_per_gpu" gorm:"not null"`
	MemoryPerGpu   int     `json:"memory_per_gpu" gorm:"not null"` // Unit `GB`
	VolumeTotal    int     `json:"volume_total" gorm:"not null"`   // Unit `GB`
	Price          float64 `json:"price" gorm:"not null"`          // 1 GPU Per Hour
	PriceVolume    float64 `json:"price_volume" gorm:"not null"`   // 1GB Per Hour

	GpuUsed    int `json:"gpu_used" gorm:"not null,default:0"`
	VolumeUsed int `json:"volume_used" gorm:"not null,default:0"`

	CreatedAt time.Time      `json:"create_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
