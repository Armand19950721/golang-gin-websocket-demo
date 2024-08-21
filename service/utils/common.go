package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"service/model"
	"service/protos/Common"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

var (
	envFileName            = ".env"
	_                      = InitEnv()
	DateTimeFormatRegular  = "2006-01-02 15:04:05"
	DateTimeStaticDate     = "2006-01-02"
	DateTimeFormatFile     = "2006_01_02T15_04_05"
	DateTimeFormatTimeZone = "2006-01-02T15:04:05-07"
)

func ValidId(str string, params ...string) bool {
	nullable := validParamsParse(params)

	// PrintObj(str, "str ValidId")
	// PrintObj(nullable, "nullable ValidId")

	if nullable && str == "" {
		return true
	}

	if str == "00000000-0000-0000-0000-000000000000" {
		return false
	}

	return govalidator.IsUUID(str)
}

func tryLoadEnv(location string) bool {
	loadEnvErr := godotenv.Load(location)
	if loadEnvErr == nil {
		PrintObj("load env success")
		return true
	}
	return false
}

func Now() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02 15:04:05")
}

func InitEnv() bool {
	envLocations := []string{
		"./" + envFileName,
		"../" + envFileName,
		"../../" + envFileName,
	}

	for _, location := range envLocations {
		if tryLoadEnv(location) {
			return true
		}
		PrintObj("load env err. try another location")
	}

	PrintObj("load env err. fail")
	return false
}

func GetEnv(str string) string {
	return os.Getenv(str)
}

func GetEnvPanic(str string) string {
	val := os.Getenv(str)
	if val == "" {
		panic("env field:" + str + " is no value")
	}

	return val
}

func IsDEV() bool {
	res := GetEnv("ENV") == "dev"
	PrintObj(res, "IsDEV")
	return res
}

func GetDomain() string {
	// set domian
	domain := "https://picbot.spe3d.co/"

	if IsDEV() {
		domain = "https://dev-picbot.spe3d.co/"
	}

	PrintObj(domain, "domain")

	return domain
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

func ToString(val int) string {
	return strconv.Itoa(val)
}

type ErrorType struct {
	Code        Common.ErrorCodes
	ReturnMsg   string
	InternalMsg string
}

func GetError(err ErrorType) string {
	if err.InternalMsg != "" {
		PrintObj(err.InternalMsg, "internalMsg")
	}

	model := &Common.ErrorReply{
		Code:    err.Code,
		Message: err.ReturnMsg,
	}

	return ToJson(model)
}

func ToJson(obj interface{}) string {
	mdJson, err := json.Marshal(obj)

	if err != nil {
		fmt.Println("to json err")
		return ""
	}

	return string(mdJson)
}

func FixStringLen(str string) string {
	limit := 500

	if len(str) > limit {
		return str[0:limit] + "....more"
	}
	return str
}

func ToInt(number string) int {
	val, err := strconv.Atoi(number)
	if err != nil {
		PrintObj("ToInt", err.Error())
		return -1
	}

	return val
}

func ToIntWithError(number string) (int, error) {
	return strconv.Atoi(number)
}

func ValidString(str string, min, max int, params ...string) bool {
	nullable := validParamsParse(params)
	// PrintObj(str, "str ValidString")
	// PrintObj(GetLength(str))
	// PrintObj(nullable)

	if nullable && str == "" {
		return true
	}

	if max == -1 { //ignore max
		if GetLength(str) >= min {
			return true
		}
	} else {
		if GetLength(str) <= max && GetLength(str) >= min {
			return true
		}
	}

	return false
}

func ValidIntFromReq(numStr string, min, max int, params ...string) bool {
	nullable := validParamsParse(params)
	// PrintObj(numStr, "ValidIntFromReq numStr")
	// PrintObj(min, "min")
	// PrintObj(max, "max")
	// PrintObj(nullable)

	if nullable && numStr == "" {
		return true
	}

	num := ToInt(numStr)

	if num == -1 {
		return false
	}

	if num > max || num < min {
		return false
	}

	return true
}

func ValidEmail(str string, params ...string) bool {
	nullable := validParamsParse(params)
	// PrintObj(str, "str ValidString")
	// PrintObj(GetLength(str))
	// PrintObj(nullable)

	if nullable && str == "" {
		return true
	}

	if !govalidator.IsEmail(str) {
		return false
	}

	return true
}

func validParamsParse(params []string) bool {

	if len(params) > 0 {
		if params[0] != "" {
			return true
		}
	}

	return false
}

func GetLength(str string) int {
	return utf8.RuneCountInString(str)
}

func GetMetaDataField(md metadata.MD, name string) string {
	arr := md[name]
	if len(arr) < 1 {
		return ""
	}
	return arr[0]
}

func ExtractToken(authHeader string) string {
	// 將字串按空格分割為一個字串陣列
	parts := strings.Split(authHeader, " ")

	// 如果陣列長度不為 2，或者第一個元素不是 "Bearer"，則返回空字串
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	// 返回第二個元素，即 Token 字串
	return parts[1]
}

func GetErrorGin(err ErrorType) gin.H {
	if err.InternalMsg != "" {
		PrintObj(err.InternalMsg, "internalMsg")
	}

	model := gin.H{
		"Code":    err.Code,
		"Message": err.ReturnMsg,
	}

	return model
}

type GinResult struct {
	Code    Common.ErrorCodes
	Message string
	Data    interface{}
}

func GetGinResult(result GinResult) gin.H {
	model := gin.H{
		"Code": result.Code,
		"Data": ToJson(result.Data),
	}

	if result.Message != "" || result.Code != Common.ErrorCodes_SUCCESS {
		PrintObj(result.Message, "result.Message")

		// add log
		// CreateLog("GetGinResult error CodeState", Common.ErrorCodes_name[int32(result.Code.Number())])
		// CreateLog("GetGinResult error Message", result.Message)
	}

	PrintObj(Common.ErrorCodes_name[int32(result.Code.Number())], "CodeState")
	PrintObj(model, "GetGinResult")

	return model
}

func GetNewUUID() uuid.UUID {
	return uuid.New()
}

func MD5Hash(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func IsErrorNotFound(err error) bool {
	res := errors.Is(err, gorm.ErrRecordNotFound)
	PrintObj(res, "IsErrorNotFound")

	return res
}

func ParseUUID(str string) uuid.UUID {
	id, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil
	}

	return id
}

func IsFileInPath(path, filename string) (bool, error) {
	files, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return false, err
	}

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		if !info.IsDir() && filepath.Base(file) == filename {
			PrintObj("found file !!")
			return true, nil
		}
	}

	PrintObj("file not found~!")
	return false, nil
}

func GetNewFileName(ctx *gin.Context) (error, string) {
	// read file
	file, err := ctx.FormFile("file") //get file
	if err != nil {
		return errors.New("fail to get form file"), ""
	}

	return nil, uuid.New().String() + GetSuffix(file)
}

func GetSuffix(file *multipart.FileHeader) string {
	uploadFileNameWithSuffix := path.Base(file.Filename)
	PrintObj(file.Filename, "file.Filename")
	PrintObj(path.Ext(uploadFileNameWithSuffix), "path.Ext(uploadFileNameWithSuffix)")
	return path.Ext(uploadFileNameWithSuffix)
}

func ParseJsonWithType[T any](str string) (res T, err error) {

	if len(str) == 0 {
		PrintObj("input invalid:", "ParseJsonWithType err")
		return res, errors.New("input invalid")
	}

	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		PrintObj(err.Error(), "ParseJsonWithType err")
		return res, err
	}

	return res, nil
}

// IsValidDate 检查日期字符串是否有效
func IsValidDate(dateStr string, params ...string) bool {

	nullable := validParamsParse(params)

	if nullable && dateStr == "" {
		return true
	}

	if dateStr == "" {
		return false
	}

	_, err := time.Parse(DateTimeFormatRegular, dateStr)
	return err == nil
}

func IsValidDateForTimeZone(dateStr string, params ...string) bool {
	PrintObj(dateStr, "IsValidDateForTimeZone from")

	nullable := validParamsParse(params)

	if nullable && dateStr == "" {
		return true
	}

	if dateStr == "" {
		return false
	}

	parseDate, err := time.Parse(DateTimeFormatTimeZone, dateStr)

	if err != nil {
		PrintObj(err.Error(), "IsValidDateForTimeZone err")
		return false
	}

	PrintObj(parseDate, "IsValidDateForTimeZone to")
	return true
}

func ParseDateOrNilTimeZone(dateString string) *time.Time {
	if dateString == "" {
		return nil
	}
	t, err := time.Parse(DateTimeFormatTimeZone, dateString)
	if err != nil {
		PrintObj(err.Error(), "ParseDateOrNil err")
		return nil
	}

	return &t
}

func ParseJson(str string) interface{} {
	var result interface{}

	err := json.Unmarshal([]byte(str), &result)

	if err != nil {
		fmt.Println("parse json err")
		return ""
	}

	return result
}

func GetFileNameAndExtension(filename string) (string, string) {
	parts := strings.Split(filename, ".")

	// 如果只有一个点，说明没有扩展名
	if len(parts) == 1 {
		return parts[0], ""
	}

	// 否则，返回文件名和扩展名
	return parts[0], parts[len(parts)-1]
}

func DeleteFiles(dirPath string) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		filePath := filepath.Join(dirPath, fileInfo.Name())
		err := os.Remove(filePath)
		if err != nil {
			PrintObj(err.Error())
		}
	}

	return nil
}

func DeleteFilesAndFolder(dirPath string) error {
	err := DeleteFiles(dirPath)
	if err != nil {
		return err
	}

	err = os.RemoveAll(dirPath)
	if err != nil {
		PrintObj(err.Error())
	}

	return nil
}

func GetPrivateLocalStaticPath(loaction string) string {
	localPath := "./staticDirPrivate/"

	if loaction != "" {
		return localPath + loaction + "/"
	}

	return loaction
}

type ReplyList struct {
	Rows       interface{}
	TotalCount int64
}

type ReplyListWithSum[T any] struct {
	Rows       []T
	RowsSum    T
	TotalCount int64
}

// EnumConverter 是一个函数类型，用于将 int32 转换为特定的枚举类型
type EnumConverter[T any] func(int32) T

// ConvertStringToEnum 通用函数，将字符串转换为枚举类型
func ConvertStringToEnum[T any](s string, nameMap map[int32]string, valueMap map[string]int32, converter EnumConverter[T]) (T, error) {
	// 将字符串转换为整数
	val, err := strconv.Atoi(s)
	if err != nil {
		return *new(T), fmt.Errorf("error converting string to int: %w", err)
	}

	// 使用整数从枚举映射中查找对应的枚举名称
	enumName, ok := nameMap[int32(val)]
	if !ok {
		return *new(T), fmt.Errorf("no enum name found for value %d", val)
	}

	// 使用枚举名称查找对应的整数值
	enumValue, ok := valueMap[enumName]
	if !ok {
		return *new(T), fmt.Errorf("no enum value found for name %s", enumName)
	}

	// 使用转换函数将 int32 转换为枚举类型
	return converter(enumValue), nil
}

func RemovePicbotSocketToken(picbots []model.Picbot) []model.Picbot {
	// filter out test picbot if is prod env
	var filter []model.Picbot
	for _, picbot := range picbots {
		picbot.SocketToken = ""
		filter = append(filter, picbot)
	}
	return filter
}

// contains 檢查字符串是否存在於陣列中
func Contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func CreateFileIfNotExist(localFilePath string) error {
	//Create directory if it doesn't exist
	dirPath := filepath.Dir(localFilePath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// 檢查手機條碼載具是否符合規範
func IsValidCarrierBarcode(barcode string) bool {
	// 檢查總長度為8碼，第一碼為「/」，其餘7碼符合規定的字符
	pattern := `^/[A-Z0-9\.\-\+]{7}$`
	matched, err := regexp.MatchString(pattern, barcode)
	if err != nil {
		fmt.Println("正則表達式錯誤:", err)
		return false
	}
	return matched
}
