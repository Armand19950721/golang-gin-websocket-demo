package httpClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	// "mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	// "path/filepath"
	"service/utils"
)

func GetInvoicePrivateKeyHeader() map[string]string {
	headers := make(map[string]string)
	headers["authorization"] = "Bearer " + utils.GetEnv("INVOICE_TOKEN")
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return headers
}

func GetPicbotPrivateKeyHeader() map[string]string {
	headers := make(map[string]string)
	headers["authorization"] = "Bearer " + utils.GetEnv("PICBOT_API_TOKEN_OLD")
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return headers
}

func GetNullHeader() map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return headers
}

func SendGETRequest(url string, headers map[string]string) string {
	utils.PrintObj(url, "SendGETRequest url")

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err.Error())
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}

	res := string(responseBody)

	// utils.PrintObj(res, "SendGETRequest")

	if res == "" {
		panic("WTF")
	}

	return string(responseBody)
}

func SendHeaderGETRequest(urlTarget string, headers map[string]string) (string, error) {
	// 建立 GET 请求
	req, err := http.NewRequest("GET", urlTarget, nil)
	if err != nil {
		return "", err
	}

	// 设置头部信息
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 将响应字节切片转换为字符串
	responseString := string(responseBytes)

	return responseString, nil
}

func SendPOSTRequest(urlTarget string, formFields map[string]string, headers map[string]string) string {
	utils.PrintObj(urlTarget, "SendPOSTRequest url")

	// 创建 url.Values，用于保存经过 URL 编码的表单字段
	formData := url.Values{}

	// 将表单字段添加到 formData 中，并自动进行 URL 编码
	for key, value := range formFields {
		formData.Set(key, value)
	}

	// 建立 POST 请求
	req, err := http.NewRequest("POST", urlTarget, strings.NewReader(formData.Encode()))
	if err != nil {
		panic(err.Error())
	}

	// 設定標頭
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 發送請求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	// 解析回應內容為 JSON
	var responseJSON map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseJSON)
	if err != nil {
		panic(err.Error())
	}

	// 將解析後的 JSON 轉換為 JSON 字串
	responseJSONString, err := json.Marshal(responseJSON)
	if err != nil {
		panic(err.Error())
	}

	res := string(responseJSONString)

	// utils.PrintObj(res, "SendPOSTRequest")

	if res == "" {
		panic("WTF")
	}

	return res
}

func SendPOSTRequestReturnString(urlTarget string, formFields map[string]string, headers map[string]string) string {
	utils.PrintObj(urlTarget, "SendPOSTRequest url")

	// 创建 url.Values，用于保存经过 URL 编码的表单字段
	formData := url.Values{}

	// 将表单字段添加到 formData 中，并自动进行 URL 编码
	for key, value := range formFields {
		formData.Set(key, value)
	}

	// 建立 POST 请求
	req, err := http.NewRequest("POST", urlTarget, strings.NewReader(formData.Encode()))
	if err != nil {
		panic(err.Error())
	}

	// 設定標頭
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 發送請求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	// 解析回應內容為 JSON
	var responseJSON string
	err = json.NewDecoder(resp.Body).Decode(&responseJSON)
	if err != nil {
		panic(err.Error())
	}

	// 將解析後的 JSON 轉換為 JSON 字串
	responseJSONString, err := json.Marshal(responseJSON)
	if err != nil {
		panic(err.Error())
	}

	res := string(responseJSONString)

	// utils.PrintObj(res, "SendPOSTRequest")

	if res == "" {
		panic("WTF")
	}

	return res
}

func SendFileUploadRequest(url string, filePath string, formFields map[string]string, headers map[string]string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fmt.Println("SendFileUploadRequest => " + filePath)

	// Add file to multipart form
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)

	// Add other form fields
	for key, value := range formFields {
		_ = writer.WriteField(key, value)
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Add("Content-Type", writer.FormDataContentType())
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var responseMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		return "", err
	}

	return utils.ToJson(responseMap), nil
}

func PutUploadFileRequest(preSignedURL string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取文件的大小
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	// 创建一个新的 HTTP PUT 请求
	req, err := http.NewRequest("PUT", preSignedURL, file)
	if err != nil {
		return err
	}

	// 设置请求头部，确保与预签名 URL 的 ContentType 一致
	req.Header.Set("Content-Type", "image/png")
	// 设置内容长度
	req.Header.Set("Content-Length", fmt.Sprintf("%d", fileSize))

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file: status code %d", resp.StatusCode)
	}

	fmt.Println("File uploaded successfully")
	return nil
}

func DownloadFile(url, filePath string) error {
	// 發送 HTTP GET 請求來下載檔案
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("無法連線到指定的 URL：%v", err)
	}
	defer response.Body.Close()

	// 創建檔案來保存下載的內容
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("無法創建檔案：%v", err)
	}
	defer file.Close()

	// 將下載的內容寫入到檔案中
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("無法寫入檔案：%v", err)
	}

	fmt.Printf("已成功下載檔案到：%s\n", filePath)
	return nil
}

func DownloadFileWithHeader(url, filePath string, headers map[string]string) error {
	utils.PrintObj(url, "DownloadFileWithHeader url")
	// 創建一個 HTTP 請求
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("無法創建 HTTP 請求：%v", err)
	}

	// 添加提供的標頭到請求中
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	// 發送請求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("無法連線到指定的 URL：%v", err)
	}
	defer response.Body.Close()

	// 創建檔案來保存下載的內容
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("無法創建檔案：%v", err)
	}
	defer file.Close()

	// 將下載的內容寫入到檔案中
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("無法寫入檔案：%v", err)
	}

	fmt.Printf("已成功下載檔案到：%s\n", filePath)
	return nil
}
