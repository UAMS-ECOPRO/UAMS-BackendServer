package models

import (
	"context"
	"fmt"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
	"time"
)

type UHF struct {
	GormModel
	UHFSerialNumber string `gorm:"type:varchar(256);unique;not null" json:"uhf_serial_number"`
	Description     string `json:"description"`
	GatewayID       string `gorm:"type:varchar(256);" json:"gateway_id"`
	ConnectState    string `json:"connect_state"`
	AreaId          string `json:"area_id"`
	UHFAddress      string `json:"uhf_address"`
	ActiveState     string `json:"active_state"`
}

// Struct defines HTTP request payload for openning doorlock
type UHFCmd struct {
	ID       string `json:"id"`
	State    string `json:"state"`
	Duration string `json:"duration"`
}

// Struct defines HTTP request payload for deleting doorlock
type UHFDelete struct {
	ID string `json:"id" binding:"required"`
}

// Struct defines HTTP request payload for getting doorlock status
type UHFStatus struct {
	ID              string `json:"id"`
	GatewayID       string `json:"gatewayId"`
	DoorlockAddress string `json:"doorlockAddress"`
	ConnectState    string `json:"connectState"`
	DoorState       string `json:"doorState"`
	LockState       string `json:"lockState"`
}

type UHFSvc struct {
	db *gorm.DB
}
type UHF_important_info struct {
	UHFAddress string `json:"uhf_address"`
}
type UHF_list struct {
	uhfs []UHF_important_info
}

func NewUHFSvc(db *gorm.DB) *UHFSvc {
	return &UHFSvc{
		db: db,
	}
}

func (uhfs *UHFSvc) FindAllUHF(ctx context.Context) (dlList []UHF, err error) {
	result := uhfs.db.Find(&dlList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlList, nil
}

func (uhfs *UHFSvc) FindUHFByID(ctx context.Context, id string) (dl *UHF, err error) {
	result := uhfs.db.First(&dl, id)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dl, nil
}

func (uhfs *UHFSvc) FindUHFByAddress(ctx context.Context, address string, gwID string) (dl *UHF, err error) {
	var cnt int64
	result := uhfs.db.Where("uhf_address = ? AND gateway_id = ?", address, gwID).Find(&dl).Count(&cnt)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if cnt <= 0 {
		return nil, fmt.Errorf("find no records")
	}

	return dl, nil
}

func (uhfs *UHFSvc) CreateUHF(ctx context.Context, dl *UHF) (*UHF, error) {
	if err := uhfs.db.Create(&dl).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dl, nil
}

func (uhfs *UHFSvc) UpdateUHF(ctx context.Context, dl *UHF) (bool, error) {
	result := uhfs.db.Model(&dl).Where("uhf_address = ? AND gateway_id = ?", dl.UHFAddress, dl.GatewayID).Updates(dl)
	return utils.ReturnBoolStateFromResult(result)
}

func (uhfs *UHFSvc) UpdateUHFByAddress(ctx context.Context, dl *UHF) (bool, error) {
	result := uhfs.db.Model(&dl).Where("gateway_id = ? AND doorlock_address = ?", dl.GatewayID, dl.UHFAddress).Updates(dl)
	return utils.ReturnBoolStateFromResult(result)
}

func (uhfs *UHFSvc) UpdateUHFState(ctx context.Context, dl *DoorlockCmd) (bool, error) {
	result := uhfs.db.Model(&Doorlock{}).Where("id = ?", dl.ID).Update("last_open_time", time.Now().UnixMilli())
	return utils.ReturnBoolStateFromResult(result)
}

func (uhfs *UHFSvc) DeleteUHF(ctx context.Context, id string) (bool, error) {
	result := uhfs.db.Unscoped().Where("id = ?", id).Delete(&UHF{})
	return utils.ReturnBoolStateFromResult(result)
}

func (uhfs *UHFSvc) DeleteUHFByAddress(ctx context.Context, dl *UHF) (bool, error) {
	result := uhfs.db.Unscoped().Where("gateway_id = ? AND doorlock_address = ?", dl.GatewayID, dl.UHFAddress).Delete(&Doorlock{})
	return utils.ReturnBoolStateFromResult(result)
}

func (uhfs *UHFSvc) UpdateUHFStatus(ctx context.Context, dl *DoorlockStatus) (bool, error) {
	result := uhfs.db.Model(&Doorlock{}).Where("id = ?", dl.ID).Updates(Doorlock{DoorState: dl.DoorState, LockState: dl.LockState})
	return utils.ReturnBoolStateFromResult(result)
}

func (uhfs *UHFSvc) GetUHFStatusByID(ctx context.Context, id string) (dl *DoorlockStatus, err error) {
	var cnt int64
	result := uhfs.db.Model(&Doorlock{}).Where("id = ?", id).Find(&dl).Count(&cnt)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if cnt <= 0 {
		return nil, fmt.Errorf("find no records")
	}

	return dl, nil
}

func (uhfs *UHFSvc) UpdateUHFStateCmd(ctx context.Context, dl *DoorlockCmd) (bool, error) {
	result := uhfs.db.Model(&Doorlock{}).Where("id = ?", dl.ID).Update("lock_state", dl.State)
	return utils.ReturnBoolStateFromResult(result)
}

func (uhfs *UHFSvc) FindAllUHFByRoomID(ctx context.Context, roomId string) (dl []*UHF, err error) {
	var cnt int64
	result := uhfs.db.Model(&UHF{}).Where("room_id = ?", roomId).Find(&dl).Count(&cnt)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if cnt <= 0 {
		return nil, fmt.Errorf("find no doorlock records by roomId")
	}

	return dl, nil
}

func (uhfs *UHFSvc) FindAllUHFByGatewayID(ctx context.Context, gwId string) (dlList []UHF, err error) {
	result := uhfs.db.Where("gateway_id = ?", gwId).Find(&dlList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlList, nil
}
