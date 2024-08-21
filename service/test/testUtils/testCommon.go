package testUtils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"service/utils"
	"time"
)

// 封裝生成帶查詢參數的 URL 的函數
func BuildURLWithQueryParameters(baseURL, route string, params map[string]string) (string, error) {
	// 解析基本 URL
	u, err := url.Parse(baseURL + route)
	if err != nil {
		return "", fmt.Errorf("error parsing base URL: %w", err)
	}

	// 添加查詢參數
	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	// 返回編碼後的 URL 字符串
	return u.String(), nil
}

func GetRandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	if min >= max {
		panic("最小值必须小于最大值")
	}

	return min + rand.Intn(max-min)
}

func GetRandomNumString() string {
	rand.Seed(time.Now().UnixNano())
	return utils.ToString(1 + rand.Intn(9999))
}

func CheckSuccess(jsonData string, shouldSuccess bool) string {
	res := false
	data := ""

	code := GetValueFromJSON(jsonData, "Code")

	if code == "" {
		panic("WTF")
	}

	if code == "10000" {
		res = true
		data = GetValueFromJSON(jsonData, "Data")
		utils.PrintObj(data, "CheckSuccess data")
	} else {
		utils.PrintObj(code, "code")
	}

	if !shouldSuccess { // 如果預測為失敗 而且res也為 false 那就算成功(返回true)
		res = !res
	}

	if !res {
		panic("didnt pass")
	}

	return data
}

func GetValueFromJSON(jsonStr string, fieldName string) string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	value, found := data[fieldName]
	if !found {
		fmt.Println("Field " + fieldName + " not found in JSON")
		return ""
	}

	// Convert value to string
	valueStr, ok := value.(string)
	if !ok {
		valueBytes, err := json.Marshal(value)
		if err != nil {
			fmt.Println(err.Error())
			return ""
		}
		valueStr = string(valueBytes)
	}

	return valueStr
}

func PrintRes(data string, params ...string) {
	// print
	key := ""

	if len(params) == 1 {
		if params[0] != "" {
			key = params[0]
			fmt.Println("=== " + key + " ===")
		}
	}

	if data != "" {
		fmt.Println(string(data))
	}
}
