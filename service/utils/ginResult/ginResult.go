package ginResult

import (
	"encoding/json"
	"fmt"
	"service/dao"
	"service/protos/Common"

	"github.com/gin-gonic/gin"
)

type CommonResult struct {
	Code    Common.ErrorCodes
	Message string
	Data    interface{}
}
type CommonResultWithType[T any] struct {
	Code    Common.ErrorCodes
	Message string
	Data    T
}

func GetGinResult(result CommonResult) gin.H {
	model := gin.H{
		"Code": result.Code,
		"Data": ToJson(result.Data),
	}

	if result.Message != "" {
		PrintObj(result.Message, "result.Message")
		dao.CreateLog("GetGinResult error Message", result.Message)
	}

	if result.Code != Common.ErrorCodes_SUCCESS {
		dao.CreateLog("GetGinResult error CodeState", Common.ErrorCodes_name[int32(result.Code.Number())])
	}

	PrintObj(Common.ErrorCodes_name[int32(result.Code.Number())], "CodeState")
	PrintObj(model, "GetGinResult")

	return model
}

func GetGinResultWithType[T any](result CommonResultWithType[T]) gin.H {
	model := gin.H{
		"Code": result.Code,
		"Data": ToJson(result.Data),
	}

	if result.Message != "" {
		PrintObj(result.Message, "result.Message")
		dao.CreateLog("GetGinResult error Message", result.Message)
	}

	if result.Code != Common.ErrorCodes_SUCCESS {
		dao.CreateLog("GetGinResult error CodeState", Common.ErrorCodes_name[int32(result.Code.Number())])
	}

	PrintObj(Common.ErrorCodes_name[int32(result.Code.Number())], "CodeState")
	PrintObj(model, "GetGinResult")

	return model
}

func ToJson(obj interface{}) string {
	mdJson, err := json.Marshal(obj)

	if err != nil {
		fmt.Println("to json err")
		return ""
	}

	return string(mdJson)
}

func PrintObj(obj interface{}, params ...string) {
	// print
	json, _ := json.Marshal(obj)
	key := ""

	if len(params) == 1 {
		if params[0] != "" {
			key = params[0]
			fmt.Println("=== " + key + " ===")
		}
	}

	if obj != "" {
		fmt.Println(string(json))
	}
}
