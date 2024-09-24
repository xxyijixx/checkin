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
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
