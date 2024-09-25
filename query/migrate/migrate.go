package main

import (
	"checkin/config"
	"checkin/query/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	var err error
	dsn := config.EnvConfig.GetDSN()
	db, err := gorm.Open(
		mysql.New(
			mysql.Config{
				DSN: dsn,
			}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	if err != nil {
		panic(fmt.Errorf("db connection failed: %v", err))
	}
	err = db.AutoMigrate(&model.UserCheckinMachine{}, &model.UserCheckinMachineInfo{}, &model.UserCheckinMachineRecord{})
	if err != nil {
		panic(fmt.Errorf("db migrate failed: %v", err))
	}
}
