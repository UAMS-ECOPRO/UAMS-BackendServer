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

func (dlsls *UHFStatusLogSvc) GetAllUHFStatusLogs(ctx context.Context) (dlslList []UHFStatusLog, err error) {
	result := dlsls.db.Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *UHFStatusLogSvc) GetUHFStatusLogByDoorID(ctx context.Context, doorId string) (dlslList []UHFStatusLog, err error) {
	result := dlsls.db.Where("uhf_address = ?", doorId).Find(&dlslList)
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

func (dlsls *UHFStatusLogSvc) GetUHFStatusLogInTimeRange(from string, to string, gateway_id string, uhf_address string) (dlslList *[]UHFStatusLog, err error) {
	result := dlsls.db.Where("created_at >= ? AND created_at <= ? AND gateway_id = ? AND uhf_address = ?", from, to, gateway_id, uhf_address).Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *UHFStatusLogSvc) DeleteUHFStatusLogInTimeRange(from string, to string) (bool, error) {
	result := dlsls.db.Unscoped().Where("created_at >= ? AND created_at <= ?", from, to).Delete(&UHFStatusLog{})
	return utils.ReturnBoolStateFromResult(result)
}

func (dlsls *UHFStatusLogSvc) DeleteUHFStatusLogUHFID(doorId string) (bool, error) {
	result := dlsls.db.Unscoped().Where("door_id = ?", doorId).Delete(&UHFStatusLog{})
	return utils.ReturnBoolStateFromResult(result)
}
