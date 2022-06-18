package models

import (
	"context"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
	"time"
)

type UserAccess struct {
	ID     uint      `gorm:"primarykey;" json:"id"`
	UserID string    `gorm:"type:varchar(256);" json:"user_id"`
	Random string    `gorm:"type:varchar(50);" json:"random"`
	Group  string    `gorm:"type:varchar(256);" json:"group"`
	AreaID string    `gorm:"type:varchar(256);" json:"area_id"`
	Time   time.Time `swaggerignore:"true" json:"time"`
}

type UserAccessSvc struct {
	db *gorm.DB
}

func NewUserAccessSvc(db *gorm.DB) *UserAccessSvc {
	return &UserAccessSvc{
		db: db,
	}
}

func (gwns *UserAccessSvc) CreateUserAccess(ctx context.Context, user_acesses *UserAccess) (*UserAccess, error) {
	if err := gwns.db.Create(&user_acesses).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return user_acesses, nil
}

func (gwns *UserAccessSvc) FindAllUserAccess(ctx context.Context) (user_acesses []UserAccess, err error) {
	result := gwns.db.Find(&user_acesses)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return user_acesses, err
}

func (gwns *UserAccessSvc) FindAllUserAccessByUserID(ctx context.Context, id string) (user_acesses []UserAccess, err error) {
	result := gwns.db.Where("user_id = ?", id).Find(&user_acesses)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return user_acesses, err
}

func (gwns *UserAccessSvc) FindAllUserAccessByUserIDAndAreaID(ctx context.Context, id string, area_id string) (user_acesses []UserAccess, err error) {
	result := gwns.db.Where("user_id = ? AND area_id = ?", id, area_id).Find(&user_acesses)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return user_acesses, err
}

func (ls *UserAccessSvc) FindUserAccessesByUserIDAndTimeRange(user_id string, from string, to string) (user_acesses *[]UserAccess, err error) {
	result := ls.db.Where("user_id = ? AND time >= ? AND time <= ?", user_id, from, to).Find(&user_acesses)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return user_acesses, nil
}

func (gwns *UserAccessSvc) FindAllUserAccessByUserIDAndAreaIDinTimeRange(ctx context.Context, id string, area_id string, from string, to string) (user_acesses []UserAccess, err error) {
	result := gwns.db.Where("user_id = ? AND area_id = ? AND time >= ? AND time <= ?", id, area_id, from, to).Find(&user_acesses)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return user_acesses, err
}

func (gwns *UserAccessSvc) FindAllUserAccessByAreaID(ctx context.Context, area_id string) (user_acesses []UserAccess, err error) {
	result := gwns.db.Where("area_id = ?", area_id).Find(&user_acesses)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return user_acesses, err
}

func (gwns *UserAccessSvc) FindAllUserAccessByAreaIDAndTimeRange(ctx context.Context, area_id string, from string, to string) (user_acesses []UserAccess, err error) {
	result := gwns.db.Where("area_id = ? AND time >= ? AND time <= ?", area_id, from, to).Find(&user_acesses)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return user_acesses, err
}

func (gwns *UserAccessSvc) FindAllUserAccessTimeRange(ctx context.Context, from string, to string) (user_acesses []UserAccess, err error) {
	result := gwns.db.Where("time >= ? AND time <= ?", from, to).Find(&user_acesses)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return user_acesses, err
}
