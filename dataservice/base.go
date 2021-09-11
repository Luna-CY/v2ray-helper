package dataservice

import (
	"errors"
	"github.com/Luna-CY/v2ray-helper/database"
	"gorm.io/gorm"
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

// UpdateByCondition 通过条件批量更新
func (b *baseDataService) UpdateByCondition(model interface{}, updates interface{}, condition interface{}, params ...interface{}) error {
	query := database.GetMainDb().Model(model)
	if nil != condition {
		query.Where(condition, params...)
	}

	return query.UpdateColumns(updates).Error
}

// Create 统一创建数据
func (b *baseDataService) Create(dest interface{}) error {
	return database.GetMainDb().Model(dest).Create(dest).Error
}

// TakeByCondition 通过条件查询单条数据
func (b *baseDataService) TakeByCondition(dest interface{}, order interface{}, condition interface{}, params ...interface{}) error {
	query := database.GetMainDb().Model(dest).Where(condition, params...)
	if nil != order {
		query.Order(order)
	}
	query.Where("deleted = 0")

	return query.Take(dest).Error
}

// TakeByConditionIgnoreDeleted 查询单条数据并忽略删除标志位
func (b *baseDataService) TakeByConditionIgnoreDeleted(dest interface{}, order interface{}, condition interface{}, params ...interface{}) error {
	query := database.GetMainDb().Model(dest).Where(condition, params...)
	if nil != order {
		query.Order(order)
	}

	return query.Take(dest).Error
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

// FindByConditionWithLimit 分页查询数据
func (b *baseDataService) FindByConditionWithLimit(page, limit int, dest interface{}, order interface{}, condition interface{}, params ...interface{}) error {
	if 0 >= page {
		return errors.New("查询页数不能小于1")
	}

	if 0 >= limit {
		return errors.New("查询条数不能小于1")
	}

	if nil == order {
		order = "id desc"
	}

	query := database.GetMainDb().Model(dest).Offset((page - 1) * limit).Limit(limit).Order(order)
	if nil != condition {
		query.Where(condition, params...)
	}
	query.Where("deleted = 0")

	return query.Find(dest).Error
}

// FindByConditionWithPagination 以分页查询列表的场景进行查询
// 该方法在计算分页的offset时会将 limit -1来进行计算
func (b *baseDataService) FindByConditionWithPagination(page, limit int, dest interface{}, order interface{}, condition interface{}, params ...interface{}) error {
	if 0 >= page {
		return errors.New("查询页数不能小于1")
	}

	if 0 >= limit {
		return errors.New("查询条数不能小于1")
	}

	if nil == order {
		order = "id desc"
	}

	// 在分页的情况下，limit会多查询一个，但是分页时计算offset需要减1计算
	query := database.GetMainDb().Model(dest).Offset((page - 1) * (limit - 1)).Limit(limit).Order(order)
	if nil != condition {
		query.Where(condition, params...)
	}
	query.Where("deleted = 0")

	return query.Find(dest).Error
}

// FindInBatches 批量处理
func (b *baseDataService) FindInBatches(dest interface{}, call func(tx *gorm.DB, batch int) error, batchSize int, order interface{}, condition interface{}, params ...interface{}) *gorm.DB {
	if nil == order {
		order = "id desc"
	}

	query := database.GetMainDb().Model(dest)
	if nil != condition {
		query.Where(condition, params...)
	}
	query.Where("deleted = 0")

	return query.FindInBatches(dest, batchSize, func(tx *gorm.DB, batch int) error {
		return call(tx, batch)
	})
}

// CountByCondition 统一统计
func (b *baseDataService) CountByCondition(model interface{}, condition interface{}, params ...interface{}) (int64, error) {
	query := database.GetMainDb().Model(model)
	if nil != condition {
		query.Where(condition, params...)
	}

	var count int64
	err := query.Count(&count).Error

	return count, err
}
