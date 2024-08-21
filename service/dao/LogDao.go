package dao

import (
	"service/model"
	"service/utils"
	"service/utils/postgresClient"
)

func CreateLogJson(key string, value interface{}) {
	valueJson := utils.ToJson(value)

	res := postgresClient.DatabaseManager.Create(&model.Log{
		Key:   key,
		Value: valueJson,
	})

	if res.Error != nil {
		utils.PrintObj(res.Error.Error(), "SaveLog error")
	}
}

func CreateLog(key string, value string) {
	utils.PrintObj([]string{key, value}, "CreateLog")

	res := postgresClient.DatabaseManager.Create(&model.Log{
		Key:   key,
		Value: value,
	})

	if res.Error != nil {
		utils.PrintObj(res.Error.Error(), "SaveLog error")
	}
}
