package dataservice

import (
	"github.com/Luna-CY/v2ray-helper/common/database"
)

var base = new(baseDataService)

func GetBaseService() *baseDataService {
	return base
}

// baseDataService 基础服务定义
type baseDataService struct{}

// GetById 通过id查询数据
// @return error 未找到数据时将返回 gorm.ErrRecordNotFound 类型错误
func (b *baseDataService) GetById(id uint, dest interface{}) error {
	return database.GetMainDb().Model(dest).Where("id = ?", id).Take(dest).Error
}

// UpdateById 通过id更新数据
// 统一数据更新与缓存更新等操作
func (b *baseDataService) UpdateById(id uint, dest interface{}) error {
	return database.GetMainDb().Model(dest).Where("id = ?", id).Updates(dest).Error
}

// Create 统一创建数据
func (b *baseDataService) Create(dest interface{}) error {
	return database.GetMainDb().Model(dest).Create(dest).Error
}

// FindByCondition 通过条件查询全部数据
func (b *baseDataService) FindByCondition(dest interface{}, order interface{}, condition interface{}, params ...interface{}) error {
	if nil == order {
		order = "id desc"
	}

	query := database.GetMainDb().Model(dest).Order(order)
	if nil != condition {
		query.Where(condition, params...)
	}
	query.Where("deleted = 0")

	return query.Find(dest).Error
}
