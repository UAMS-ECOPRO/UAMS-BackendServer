package models

import (
	"context"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type SystemLog struct {
	GormModel
	GatewayID string `json:"gatewayId"`
	LogType   string `json:"-"` // info, warn, debug, error, fatal
	Content   string `json:"content"`
}

type SystemLogSvc struct {
	db *gorm.DB
}

func NewSystemLogSvc(db *gorm.DB) *SystemLogSvc {
	return &SystemLogSvc{
		db: db,
	}
}

func (ls *SystemLogSvc) FindAllGatewayLog(ctx context.Context) (glList []GatewayLog, err error) {
	result := ls.db.Find(&glList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return glList, nil
}

func (ls *SystemLogSvc) FindGatewayLogByID(ctx context.Context, id string) (gl *GatewayLog, err error) {
	result := ls.db.Preload("Doorlocks").First(&gl, id)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gl, nil
}

func (ls *SystemLogSvc) CreateSystemLog(ctx context.Context, gl *SystemLog) (*SystemLog, error) {
	if err := ls.db.Create(&gl).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gl, nil
}
