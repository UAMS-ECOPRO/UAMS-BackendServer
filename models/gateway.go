package models

import (
	"context"
	"fmt"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Gateway struct {
	GormModel
	AreaID          string      `json:"area_id"`
	GatewayID       string      `gorm:"type:varchar(256);unique;not null;" json:"gateway_id"`
	Name            string      `json:"name"`
	ConnectState    bool        `gorm:"type:bool;not null;" json:"connect_state"`
	SoftwareVersion string      `json:"software_version"`
	UHFs            []UHF       `gorm:"foreignKey:GatewayID;references:GatewayID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"uhfs"`
	GwNetworks      []GwNetwork `gorm:"foreignKey:GatewayID;references:GatewayID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"gw_networks"`
}

// Struct defines HTTP request payload for deleting gateway
type DeleteGateway struct {
	GatewayID string `json:"gateway_id" binding:"required"`
}

type UpdateGateway struct {
	GatewayID string `json:"gateway_id" binding:"required"`
}

type GatewayBlockCmd struct {
	BlockId string `json:"block_id" binding:"required"`
	Action  string `json:"action" binding:"required"`
}
type GatewaySvc struct {
	db *gorm.DB
}

func NewGatewaySvc(db *gorm.DB) *GatewaySvc {
	return &GatewaySvc{
		db: db,
	}
}

func (gs *GatewaySvc) FindAllGateway(ctx context.Context) (gwList []Gateway, err error) {
	result := gs.db.Preload("UHFs").Preload("GwNetworks").Find(&gwList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gwList, nil
}

func (gs *GatewaySvc) FindGatewayByID(ctx context.Context, id string) (gw *Gateway, err error) {
	result := gs.db.Preload("UHFs").Preload("GwNetworks").First(&gw, id)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gw, nil
}

func (gs *GatewaySvc) FindGatewayByMacID(ctx context.Context, id string) (gw *Gateway, err error) {
	var cnt int64
	result := gs.db.Preload("GwNetworks").Where("gateway_id = ?", id).Find(&gw).Count(&cnt)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if cnt <= 0 {
		return nil, fmt.Errorf("find no records")
	}

	return gw, nil
}

func (gs *GatewaySvc) CreateGateway(ctx context.Context, g *Gateway) (*Gateway, error) {
	if err := gs.db.Create(&g).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return g, nil
}

func (gs *GatewaySvc) UpdateGateway(ctx context.Context, g *Gateway) (bool, error) {
	var cnt int64
	gateway := gs.db.Model(&g).Where("gateway_id = ?", g.GatewayID)
	gateway.Count(&cnt)
	if cnt <= 0 {
		return false, fmt.Errorf("No gateway found")
	}
	result := gateway.Updates(g)
	return utils.ReturnBoolStateFromResult(result)
}

func (gs *GatewaySvc) DeleteGateway(ctx context.Context, gwID string) (bool, error) {
	result := gs.db.Unscoped().Where("gateway_id = ?", gwID).Delete(&Gateway{})
	return utils.ReturnBoolStateFromResult(result)
}

func (gs *GatewaySvc) AppendGatewayDoorlock(ctx context.Context, gw *Gateway, d *Doorlock) (*Gateway, error) {
	if err := gs.db.Model(&gw).Association("Doorlocks").Append(d); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gw, nil
}

func (gs *GatewaySvc) UpdateGatewayDoorlock(ctx context.Context, gw *Gateway, d *Doorlock) (*Gateway, error) {
	if err := gs.db.Model(&gw).Association("Doorlocks").Replace(d); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gw, nil
}

func (gs *GatewaySvc) DeleteGatewayDoorlock(ctx context.Context, gw *Gateway, d *Doorlock) (*Gateway, error) {
	if err := gs.db.Model(&gw).Association("Doorlocks").Delete(d); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gw, nil
}

func (gs *GatewaySvc) FindAllGatewaysByBlockID(ctx context.Context, block_id string) (gwList []string, err error) {
	if err := gs.db.Model(&Doorlock{}).Select("gateway_id").Where("block_id = ?", block_id).Group("gateway_id").Find(&gwList).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gwList, nil
}

func (gs *GatewaySvc) UpdateAllDoorlocksStateByBlockID(ctx context.Context, block_id string, state string) (bool, error) {
	result := gs.db.Model(&Doorlock{}).Where("block_id = ?", block_id).Update("lock_state", state)
	return utils.ReturnBoolStateFromResult(result)
}

func (gs *GatewaySvc) UpdateGatewayConnectState(ctx context.Context, gwId string, state bool) (bool, error) {
	if err := gs.db.Model(&Gateway{}).Where("gateway_id = ?", gwId).Update("connect_state", state).Error; err != nil {
		err = utils.HandleQueryError(err)
		return false, err
	}
	return true, nil
}