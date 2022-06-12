package models

import (
	"context"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type OperationLog struct {
	GormModel
	GatewayID string `json:"gateway_id"`
	Content   string `json:"content"` // corresponding statetype
}
type OperationLogSvc struct {
	db *gorm.DB
}

func NewOperationLogSvc(db *gorm.DB) *OperationLogSvc {
	operationLog := &OperationLogSvc{
		db: db,
	}
	return operationLog
}

func (dlsls *OperationLogSvc) GetAllOperationLogs(ctx context.Context) (dlslList []OperationLog, err error) {
	result := dlsls.db.Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *OperationLogSvc) GetOperationLogByDoorID(ctx context.Context, doorId string) (dlslList []OperationLog, err error) {
	result := dlsls.db.Where("door_id = ?", doorId).Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *OperationLogSvc) CreateOperationLog(ctx context.Context, dlsl *OperationLog) (*OperationLog, error) {
	if err := dlsls.db.Create(&dlsl).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlsl, nil
}

func (dlsls *OperationLogSvc) GetOperationLogInTimeRange(from string, to string) (dlslList *[]OperationLog, err error) {
	result := dlsls.db.Where("created_at >= ? AND created_at <= ?", from, to).Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *OperationLogSvc) DeleteOperationLogInTimeRange(from string, to string) (bool, error) {
	result := dlsls.db.Unscoped().Where("created_at >= ? AND created_at <= ?", from, to).Delete(&OperationLog{})
	return utils.ReturnBoolStateFromResult(result)
}

func (dlsls *OperationLogSvc) DeleteOperationLogByDoorID(doorId string) (bool, error) {
	result := dlsls.db.Unscoped().Where("door_id = ?", doorId).Delete(&OperationLog{})
	return utils.ReturnBoolStateFromResult(result)
}
