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
		&UserAccess{},
		&PackageAccess{},
		&SystemLog{},
		&UHFStatusLog{},
		&OperationLog{},
	)
	if err != nil {
		panic(err)
	}
}
