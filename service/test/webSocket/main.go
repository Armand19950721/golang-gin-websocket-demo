package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	// "service/model"
	"service/model"
	"service/protos/Common"
	"service/protos/PicbotEnum"
	"service/test/service"
	"service/test/testUtils"
	"service/utils/authUtils"

	"service/utils"
	"service/utils/webSocketServer"
	"time"

	"github.com/gorilla/websocket"
)

var picbotNameUse = ""
var tokenUse = ""
var CarruerNum = "/FR6HZMA"

func main() {
	headers := http.Header{}
	// parse flag
	flag.Parse()
	cmd := flag.Arg(0)
	tokens := authUtils.GetAllToken()

	if cmd == "1" {
		picbotNameUse = "PICBOT_01_TOKEN"
		tokenUse = tokens.PICBOT_01_TOKEN
	}else {
		panic("WTF you need to choise a picbot")
	}

	headers.Add("Authorization", tokenUse)

	// 使用for循环，使客户端持续尝试连接
	for {
		// 连接WebSocket服务器并传递令牌
		conn, _, err := websocket.DefaultDialer.Dial(service.ServerAddr, headers)
		if err != nil {
			log.Println("Error connecting to server:", err)
			// 尝试重连
			reconnectInterval := 5 * time.Second
			log.Printf("Attempting to reconnect in %s...\n", reconnectInterval)
			time.Sleep(reconnectInterval)
			continue
		}
		defer conn.Close()

		fmt.Println("Connected to server")

		// 启动心跳goroutine
		// go sendGetThemeGroup(conn)
		go sendOrder(conn)
		// go sendState(conn)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}

			fmt.Printf("Received message: %s\n", message)

			reply, err := utils.ParseJsonWithType[webSocketServer.SocketMessage](string(message))
			if err != nil {
				panic(err.Error())
			}

			if reply.MessageType == int32(PicbotEnum.PicbotCommandType_RECEIVE_CONFIRM) ||
				reply.MessageType == int32(PicbotEnum.PicbotCommandType_REPLY_FROM_BACKEND) {
				continue
			}

			reply.From = reply.To
			reply.To = picbotNameUse
			reply.MessageType = int32(PicbotEnum.PicbotCommandType_REPLY_FROM_PICBOT)

			replyMessage := utils.ToJson(reply)

			// Echo the message back to the client
			err = conn.WriteMessage(websocket.TextMessage, []byte(replyMessage))
			if err != nil {
				log.Println("Error writing message:", err)
				break
			}

			fmt.Printf("Reply message: %s\n", replyMessage)
		}

		log.Println("Lost connection to webSocketServer. Attempting to reconnect...")
	}

}

func sendGetThemeGroup(conn *websocket.Conn) {

	message := webSocketServer.SocketMessage{
		From:        picbotNameUse,
		To:          webSocketServer.BackendSenderNameDefault,
		State:       Common.ErrorCodes_SUCCESS,
		MessageType: int32(PicbotEnum.PicbotCommandType_GET_THEME_GROUP_INFO),
	}

	SendRequst(conn, message)
	fmt.Println("Sent GetThemeGroup")
}

func sendOrder(conn *websocket.Conn) {

	datas := []webSocketServer.AddOrderRequest{
		{ // 測試正常載具
			Order: model.Orders{
				PicbotName:     picbotNameUse,
				PaymentType:    PicbotEnum.PicbotPaymentType_EZ_CARD,
				Remark:         "test",
				State:          PicbotEnum.PicbotOrderState_SUCCESS,
				PaymentDate:    time.Now(),
				ItemPrice:      100,
				AmountPrintOut: 2,
				TotalPrice:     200,
				CarruerNumber:  CarruerNum,
			},
			OrderDetail: model.OrderDetails{
				QrCodeNumber: "1234",
				PackNameA:    "yess",
				PackNameB:    "nooo",
			},
		},
	}

	for _, data := range datas {
		message := webSocketServer.SocketMessage{
			From:        picbotNameUse,
			To:          webSocketServer.BackendSenderNameDefault,
			DataJson:    utils.ToJson(data),
			State:       Common.ErrorCodes_SUCCESS,
			MessageType: int32(PicbotEnum.PicbotCommandType_ADD_ORDER),
		}

		SendRequst(conn, message)
		fmt.Println("Sent order")
		break
	}

}

func sendState(conn *websocket.Conn) {

	heartbeatInterval := 1 * time.Second // 定义心跳间隔
	count := 0                           // 循環打不同case

	for {
		// 等待心跳间隔
		time.Sleep(heartbeatInterval)

		datas := []webSocketServer.StateChangeRequest{
			{
				Field: "NetworkSpeedUpload",
				Value: utils.ToString(testUtils.GetRandomInt(1, 100)),
			},
			{
				Field: "NetworkSpeedDownload",
				Value: utils.ToString(testUtils.GetRandomInt(1, 100)),
			},
		}

		message := webSocketServer.SocketMessage{
			From:        picbotNameUse,
			To:          webSocketServer.BackendSenderNameDefault,
			DataJson:    utils.ToJson(datas[count]),
			State:       Common.ErrorCodes_SUCCESS,
			MessageType: int32(PicbotEnum.PicbotCommandType_PICBOT_STATE),
		}

		SendRequst(conn, message)
		
		// 循環
		if count == len(datas)-1 {
			count = 0
		} else {
			count++
		}
	}
}

func SendRequst(conn *websocket.Conn, message interface{}) {
	// 构建消息
	heartbeatMessage := []byte(utils.ToJson(message))

	// 发送消息
	err := conn.WriteMessage(websocket.TextMessage, heartbeatMessage)
	if err != nil {
		log.Println("Error sending state:", err)
		return
	}
}
