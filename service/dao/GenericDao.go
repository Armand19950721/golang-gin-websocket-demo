package dao

import (
	"service/utils"
	"service/utils/postgresClient"

	"gorm.io/gorm"
)

func Create[T any](model T) (*gorm.DB, T) {
	result := postgresClient.DatabaseManager.Create(&model)

	return result, model
}

func GetList[T any](whereModel T) (tx *gorm.DB, models []T, count int64) {
	// 在 Where 和 Find 之間鏈接 Order 方法以依據 created_at 欄位降序排序
	// 預設 Limit 設為 100
	result := postgresClient.DatabaseManager.
		Where(whereModel).
		Order("create_at DESC").
		Limit(100).
		Find(&models).
		Count(&count)

	return result, models, count
}

func GetListPage[T any](whereModel T, limit, offset int) (tx *gorm.DB, models []T) {
	result := postgresClient.DatabaseManager.Where(whereModel).
		Order("create_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models)

	return result, models
}

func Query[T any](whereModel T) (tx *gorm.DB, data T) {
	result := postgresClient.DatabaseManager.Where(whereModel).First(&data)

	return result, data
}

func Update[T any](where, model T) (tx *gorm.DB, data T) {
	result := postgresClient.DatabaseManager.Where(where).Updates(&model)

	return result, model
}

func SoftDelete[T any](whereModel T) (tx *gorm.DB) {
	result := postgresClient.DatabaseManager.Where(whereModel).Delete(&whereModel)

	return result
}

// 用之前請頭腦清醒
func HardDelete[T any](whereModel T) (tx *gorm.DB) {
	result := postgresClient.DatabaseManager.Unscoped().Where(whereModel).Delete(&whereModel)

	return result
}

func CreateLogGeneric[T any](model T) {
	utils.PrintObj(model, "CreateLogGeneric")

	result := postgresClient.DatabaseManager.Create(&model)

	if result.Error != nil {
		utils.PrintObj(result.Error.Error(), "CreateLogGeneric error")
	}
}
