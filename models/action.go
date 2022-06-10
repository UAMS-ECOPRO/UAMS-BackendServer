package models

import (
	"context"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Action struct {
	GormModel
	GatewayID  string `gorm:"type:varchar(256);" json:"gateway_id"`
	UHFAddress string `gorm:"type:varchar(50);not null" json:"uhf_address"`
	EPC        string `gorm:"type:varchar(256);" json:"epc"`
}

type ActionSvc struct {
	db *gorm.DB
}

func NewActionSvc(db *gorm.DB) *ActionSvc {
	return &ActionSvc{
		db: db,
	}
}

func (gwns *ActionSvc) CreateAction(ctx context.Context, act *Action) (*Action, error) {
	if err := gwns.db.Create(&act).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return act, nil
}

func (gwns *ActionSvc) FindAction(ctx context.Context, gwId string, ifName string) (gwNet *GwNetwork, err error) {
	result := gwns.db.Where("gateway_id = ? AND interface_name = ?", gwId, ifName).Find(&gwNet)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gwNet, err
}
