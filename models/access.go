package models

import (
	"context"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Access struct {
	GormModel
	GatewayID  string `gorm:"type:varchar(256);" json:"gateway_id"`
	UHFAddress string `gorm:"type:varchar(50);not null" json:"uhf_address"`
	EPC        string `gorm:"type:varchar(256);" json:"epc"`
}

type AccessSvc struct {
	db *gorm.DB
}

func NewAccessSvc(db *gorm.DB) *AccessSvc {
	return &AccessSvc{
		db: db,
	}
}

func (gwns *AccessSvc) CreateAction(ctx context.Context, act *Access) (*Access, error) {
	if err := gwns.db.Create(&act).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return act, nil
}

func (gwns *AccessSvc) FindAllActions(ctx context.Context) (gwNet []Access, err error) {
	result := gwns.db.Find(&gwNet)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gwNet, err
}

func (gwns *AccessSvc) FindAllActionsByEPC(ctx context.Context, epc string) (gwNet []Access, err error) {
	result := gwns.db.Where("epc = ?", epc).Find(&gwNet)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gwNet, err
}
