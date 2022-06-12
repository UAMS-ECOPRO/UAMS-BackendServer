package models

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Area{},
		&Gateway{},
		&UHF{},
		&GatewayLog{},
		&GwNetwork{},
		&Action{},
		&SystemLog{},
		&UHFStatusLog{},
		&OperationLog{},
	)
	if err != nil {
		panic(err)
	}
}
