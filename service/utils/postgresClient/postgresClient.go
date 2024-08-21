package postgresClient

import (
	"fmt"
	"service/model"
	"service/utils"
	"service/utils/authUtils"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DatabaseManager = gorm.DB{}
	masterDsn       = utils.GetEnv("DB_INFO_MASTER")
)

func InitDB() bool {
	ticker := time.NewTicker(1 * time.Second)
	count := 0
	success := false

	for range ticker.C {
		success = excute()
		utils.PrintObj("InitDB retry:"+utils.ToString(count), "")

		if count == 60 || success {
			return success
		}

		count++
	}

	return true
}

func excute() bool {
	utils.PrintObj(masterDsn, "masterDsn")
	// connect master
	masterDb, err := gorm.Open(postgres.Open(masterDsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		utils.PrintObj(err.Error(), "init db err")
		return false
	}

	DatabaseManager = *masterDb
	utils.PrintObj("db master connected")

	// init uuid function
	init := masterDb.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if init.Error != nil {
		utils.PrintObj(init.Error.Error(), "gorm DB init uuid function fail")
		return false
	}

	// init table or migrate
	migrateErr := masterDb.AutoMigrate(
		&model.Log{},
		&model.Orders{},
		&model.OrderDetails{},
		&model.Picbot{},
	)

	if migrateErr != nil {
		utils.PrintObj(migrateErr, "gorm migrate error")
		return false
	}

	// init picbot test list if empty
	_, _, count := getList(model.Picbot{})

	if count == 0 {
		utils.PrintObj("creating picbots...")
		tokens := authUtils.GetAllToken()
		tokensOld := authUtils.GetAllTokenOld()
		create(model.Picbot{Name: "PICBOT_01", Remark: "", SocketToken: tokens.PICBOT_01_TOKEN, SocketTokenOld: tokensOld.PICBOT_01_TOKEN_OLD})
		create(model.Picbot{Name: "PICBOT_02", Remark: "", SocketToken: tokens.PICBOT_02_TOKEN, SocketTokenOld: tokensOld.PICBOT_02_TOKEN_OLD})
		create(model.Picbot{Name: "PICBOT_03", Remark: "", SocketToken: tokens.PICBOT_03_TOKEN, SocketTokenOld: tokensOld.PICBOT_03_TOKEN_OLD})
	}

	TestDbMasterSlave()

	return true
}

func DisConnect() bool {
	db, _ := DatabaseManager.DB()
	err := db.Close()
	fmt.Println(err)

	return err == nil
}

func TestDbMasterSlave() {
	testId := uuid.New()
	checkValue := uuid.New()

	res := DatabaseManager.Create(&model.Log{
		Id:    testId,
		Key:   "test db master slave",
		Value: checkValue.String(),
	})

	if res.Error != nil {
		panic(res.Error.Error())
	}

	time.Sleep(1 * time.Second)
	// utils.PrintObj("sleep over")

	findResult := model.Log{}
	DatabaseManager.Where(model.Log{
		Id: testId,
	}).Find(&findResult)

	if res.Error != nil {
		panic(res.Error.Error())
	}

	if findResult.Value != checkValue.String() {
		utils.PrintObj("", "test db master slave")
		panic("test fail")
	} else {
		utils.PrintObj("test success", "test db master slave")
	}
}

// just use at this file
func create[T any](model T) (*gorm.DB, T) {
	result := DatabaseManager.Create(&model)

	if result.Error != nil {
		panic(result.Error.Error())
	}

	return result, model
}

func getList[T any](whereModel T) (tx *gorm.DB, models []T, count int64) {
	result := DatabaseManager.Where(whereModel).Find(&models).Count(&count)

	if result.Error != nil {
		panic(result.Error.Error())
	}

	return result, models, count
}

func query[T any](whereModel T) (tx *gorm.DB, data T) {
	result := DatabaseManager.Where(whereModel).First(&data)

	if result.Error != nil {
		panic(result.Error.Error())
	}

	return result, data
}

func update[T any](where, model T) (tx *gorm.DB, data T) {
	result := DatabaseManager.Where(where).Updates(&model)

	if result.Error != nil {
		panic(result.Error.Error())
	}

	return result, model
}
