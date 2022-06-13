package models

import (
	"context"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
	"time"
)

type PackageAccess struct {
	ID        uint      `gorm:"primarykey;" json:"id"`
	PackageID string    `gorm:"type:varchar(256);" json:"package_id"`
	Random    string    `gorm:"type:varchar(50);" json:"random"`
	Group     string    `gorm:"type:varchar(256);" json:"group"`
	AreaID    string    `gorm:"type:varchar(256);" json:"area_id"`
	Time      time.Time `swaggerignore:"true" json:"created_at"`
}

type PackageAccessSvc struct {
	db *gorm.DB
}

func NewPackageAccessSvc(db *gorm.DB) *PackageAccessSvc {
	return &PackageAccessSvc{
		db: db,
	}
}

func (gwns *PackageAccessSvc) CreatePackageAccess(ctx context.Context, package_acesses *PackageAccess) (*PackageAccess, error) {
	if err := gwns.db.Create(&package_acesses).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return package_acesses, nil
}

func (gwns *PackageAccessSvc) FindAllPackageAccess(ctx context.Context) (package_acesses []PackageAccess, err error) {
	result := gwns.db.Find(&package_acesses)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return package_acesses, err
}

func (gwns *PackageAccessSvc) FindAllPackageAccessByPackageID(ctx context.Context, id string) (package_acesses []PackageAccess, err error) {
	result := gwns.db.Where("package_id = ?", id).Find(&package_acesses)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return package_acesses, err
}

func (ls *PackageAccessSvc) FindUserPackageByPackageIDAndTimeRange(user_id string, from string, to string) (package_acesses *[]PackageAccess, err error) {
	result := ls.db.Where("package_id = ? AND time >= ? AND time <= ?", user_id, from, to).Find(&package_acesses)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return package_acesses, nil
}
