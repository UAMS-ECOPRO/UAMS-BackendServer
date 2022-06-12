package models

import (
	"context"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type UHFLogTime struct {
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}
type UHFStatusLog struct {
	GormModel
	GatewayID  string `json:"gateway_id"`
	UHFAddress string `json:"uhf_address"`
	StateType  string `json:"state_type"`  // ConnectState, DoorState, LockState
	StateValue string `json:"state_value"` // corresponding statetype
}
type UHFStatusLogSvc struct {
	db *gorm.DB
}

func NewUHFStatusLogSvc(db *gorm.DB) *UHFStatusLogSvc {
	uhfStatusSvc := &UHFStatusLogSvc{
		db: db,
	}
	return uhfStatusSvc
}

func (dlsls *UHFStatusLogSvc) GetAllDoorlockStatusLogs(ctx context.Context) (dlslList []UHFStatusLog, err error) {
	result := dlsls.db.Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *UHFStatusLogSvc) GetDoorlockStatusLogByDoorID(ctx context.Context, doorId string) (dlslList []UHFStatusLog, err error) {
	result := dlsls.db.Where("door_id = ?", doorId).Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *UHFStatusLogSvc) CreateUHFStatusLog(ctx context.Context, dlsl *UHFStatusLog) (*UHFStatusLog, error) {
	if err := dlsls.db.Create(&dlsl).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlsl, nil
}

func (dlsls *UHFStatusLogSvc) GetDoorlockStatusLogInTimeRange(from string, to string) (dlslList *[]UHFStatusLog, err error) {
	result := dlsls.db.Where("created_at >= ? AND created_at <= ?", from, to).Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *UHFStatusLogSvc) DeleteDoorlockStatusLogInTimeRange(from string, to string) (bool, error) {
	result := dlsls.db.Unscoped().Where("created_at >= ? AND created_at <= ?", from, to).Delete(&UHFStatusLog{})
	return utils.ReturnBoolStateFromResult(result)
}

func (dlsls *UHFStatusLogSvc) DeleteDoorlockStatusLogByDoorID(doorId string) (bool, error) {
	result := dlsls.db.Unscoped().Where("door_id = ?", doorId).Delete(&UHFStatusLog{})
	return utils.ReturnBoolStateFromResult(result)
}
