package models

import (
	"time"

	"gorm.io/gorm"
)

type Orders struct {
	ID uint `json:"id" gorm:"primary_key;autoIncrement;index"`

	UserID    uint      `json:"user_id" gorm:"not null"`
	ServerID  uint      `json:"server_id" gorm:"not null"`
	StartTime time.Time `json:"start_time" gorm:"not null"`
	EndTime   time.Time `json:"end_time" gorm:"not null"`

	GpuNum   int     `json:"gpu_num" gorm:"not null"`
	DiskUsed float64 `json:"disk_used" gorm:"not null"`

	PricePerHour     float64 `json:"price_per_gpu_per_hour" gorm:"not null"`
	PriceDiskPerHour float64 `json:"price_disk_per_hour" gorm:"not null"`

	Price float64 `json:"price" gorm:"not null"`

	CreatedAt time.Time      `json:"create_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
