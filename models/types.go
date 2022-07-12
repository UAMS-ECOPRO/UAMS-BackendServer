// Package models provides business entity models and related business logics for the app
package models

import (
	"time"
)

type GormModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `swaggerignore:"true" json:"created_at"`
	UpdatedAt time.Time `swaggerignore:"true" json:"updated_at"`
}

type DeleteID struct {
	ID uint `json:"id"`
}

// Struct defines all services for our IoC
type ServiceOptions struct {
	GatewaySvc       *GatewaySvc
	AreaSvc          *AreaSvc
	LogSvc           *LogSvc
	UHFStatusLogSvc  *UHFStatusLogSvc
	UHFSvc           *UHFSvc
	UserAccessSvc    *UserAccessSvc
	PackageAccessSvc *PackageAccessSvc
	SystemLogSvc     *SystemLogSvc
	OperationLogSvc  *OperationLogSvc
}
