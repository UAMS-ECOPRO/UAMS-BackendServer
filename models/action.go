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

func (gwns *ActionSvc) FindAllActions(ctx context.Context) (gwNet []Action, err error) {
	result := gwns.db.Find(&gwNet)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gwNet, err
}

func (gwns *ActionSvc) FindAllActionsByEPC(ctx context.Context, epc string) (gwNet []Action, err error) {
	result := gwns.db.Where("epc = ?", epc).Find(&gwNet)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gwNet, err
}
