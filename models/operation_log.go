package models

import (
	"context"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type OperationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Time      time.Time `swaggerignore:"true" json:"time"`
	GatewayID string    `json:"gateway_id"`
	Content   string    `json:"content"` // corresponding statetype
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

func (dlsls *OperationLogSvc) GetOperationLogByGatewayID(ctx context.Context, doorId string) (dlslList []OperationLog, err error) {
	result := dlsls.db.Where("gateway_id = ?", doorId).Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *OperationLogSvc) FindOperationLogsByGatewayIDAndTime(gateway_id string, from string, to string) (dlslList *[]OperationLog, err error) {
	result := dlsls.db.Where("gateway_id = ? AND time >= ? AND time <= ?", gateway_id, from, to).Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *OperationLogSvc) FindOperationLogsByTime(from string, to string) (dlslList *[]OperationLog, err error) {
	result := dlsls.db.Where("time >= ? AND time <= ?", from, to).Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}
