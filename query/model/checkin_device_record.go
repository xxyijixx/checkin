package model

import (
	"checkin/config"
	"time"
)

var TableNameCheckinDeviceRecord = config.EnvConfig.DB_PREFIX + "checkin_device_records"

type CheckinDeviceRecord struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Sn         string    `gorm:"comment:设备序列号; size:255" json:"sn"`
	Mode       int       `json:"mode"`
	Inout      int       `gorm:"comment:0 主机 1 子机" json:"inout"`
	Event      int       `gorm:"comment:自定义动作" json:"event"`
	Enrollid   int       `json:"enrollid"`
	ReportTime time.Time `gorm:"comment:上报时间" json:"report_time"`
	CreatedAt  time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time `gorm:"comment:更新时间" json:"updated_at"`
}

func (*CheckinDeviceRecord) TableName() string {
	return TableNameCheckinDeviceRecord
}
