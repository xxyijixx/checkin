package model

import (
	"time"

	"gorm.io/gorm"
)

type UserCheckinMachineInfo struct {
	ID        int            `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Sn        string         `gorm:"size:255" json:"sn"`
	Enrollid  int            `json:"enrollid"`
	Name      string         `gorm:"size:255" json:"name"`
	Backupnum int            `json:"backupnum"`
	Admin     int            `json:"admin"`
	Record    string         `json:"record"`
	Status    int            `gorm:"comment:-1 未登记 0 禁用 1 启用; default:-1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
