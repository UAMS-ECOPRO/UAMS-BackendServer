package models

import (
	"context"
	"time"

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
	ID         uint      `gorm:"primaryKey" json:"id"`
	Time       time.Time `swaggerignore:"true" json:"time"`
	GatewayID  string    `json:"gateway_id"`
	UHFAddress string    `json:"uhf_address"`
	StateType  string    `json:"state_type"`  // ConnectState, DoorState, LockState
	StateValue string    `json:"state_value"` // corresponding statetype
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

func (gs *UHFStatusLogSvc) GetUHFStatusLogByID(ctx context.Context, id string) (uhf_log *UHFStatusLog, err error) {
	result := gs.db.First(&uhf_log, id)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return uhf_log, nil
}

func (dlsls *UHFStatusLogSvc) GetUHFStatusLogByUHFAddress(ctx context.Context, uhf_address string, gateway_id string) (dlslList []UHFStatusLog, err error) {
	result := dlsls.db.Where("uhf_address = ? AND gateway_id <= ?", uhf_address, gateway_id).Find(&dlslList)
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

func (dlsls *UHFStatusLogSvc) GetUHFStatusLogBYGatewayIDAndUHFAddressInTimeRange(from string, to string, gateway_id string, uhf_address string) (dlslList *[]UHFStatusLog, err error) {
	result := dlsls.db.Where("time >= ? AND time <= ? AND gateway_id = ? AND uhf_address = ?", from, to, gateway_id, uhf_address).Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *UHFStatusLogSvc) GetUHFStatusLogInTimeRange(from string, to string) (dlslList *[]UHFStatusLog, err error) {
	result := dlsls.db.Where("time >= ? AND time <= ?", from, to).Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *UHFStatusLogSvc) DeleteUHFLogInTimeRange(from string, to string) (bool, error) {
	result := dlsls.db.Unscoped().Where("time >= ? AND time <= ?", from, to).Delete(&UHFStatusLog{})
	return utils.ReturnBoolStateFromResult(result)
}
