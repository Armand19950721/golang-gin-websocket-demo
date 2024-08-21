package service

import (
	"fmt"
	"service/dao"
	"service/model"
	"service/protos/Common"
	"service/protos/PicbotEnum"
	"service/utils"
)

type PicbotAlive struct {
	PicbotName     string
	IsAlive        bool
	ConnectTime    string
	DisconnectTime string
}

type PicbotState struct {
	PicbotName           string
	PicbotBootTime       string
	Temperature          int
	RemainingPages       int
	TakePictureCount     int
	PrintCount           int
	PeopleCount          int
	NetworkSpeedUpload   int
	NetworkSpeedDownload int
}

type PicbotProductStaticDayDisplay struct {
	StaticDate         string
	PicbotName         string
	PeopleCount        int32
	TakePictureCount   int32
	PrintCount         int32
	PaymentCount       int32
	EzCardTotal        int32
	MobileTotalPayment int32
	TotalIncome        int32
	ConversionRate     string
}

func ConvertInt32ToPicbotCommandType(v int32) (PicbotEnum.PicbotCommandType, error) {
	// 使用整数从枚举映射中查找对应的枚举名称
	enumName, ok := PicbotEnum.PicbotCommandType_name[v]
	if !ok {
		return 0, fmt.Errorf("no enum name found for value %d", v)
	}

	// 使用枚举名称查找对应的整数值，并将其转换为枚举实例
	enumValue, ok := PicbotEnum.PicbotCommandType_value[enumName]
	if !ok {
		return 0, fmt.Errorf("no enum value found for name %s", enumName)
	}

	return PicbotEnum.PicbotCommandType(enumValue), nil
}

func GetPicbotInfoForPicbot(picbotName string) (code Common.ErrorCodes, data model.Picbot) {
	utils.PrintObj(picbotName, "GetPicbot")

	// check param
	if !utils.ValidString(picbotName, 1, 100) {
		code = Common.ErrorCodes_INVAILD_PARAM
		return
	}

	// query
	res, data := dao.Query(model.Picbot{Name: picbotName})
	if res.Error != nil {
		code = Common.ErrorCodes_DB_ERROR
		return
	}

	// return
	data = model.Picbot{
		Name:         data.Name,
		Remark:       data.Remark,
		ThemeGroupId: data.ThemeGroupId,
	}
	code = Common.ErrorCodes_SUCCESS

	return
}
