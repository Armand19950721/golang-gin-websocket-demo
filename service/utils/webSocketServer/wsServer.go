package webSocketServer

import (
	"fmt"
	"net/http"
	"service/dao"
	"service/protos/Common"
	"service/protos/PicbotEnum"
	"service/service"
	"service/utils"
	"service/utils/ginResult"

	"github.com/gorilla/websocket"
)

var BackendSenderNameDefault = "Backend"
var PicbotTokenMap = map[string]string{}              // 键是picbot token ，值是PICBOT编号
var WebSocketConns = make(map[string]*websocket.Conn) // WebSocket连接列表 key:picbot name ,value: conn

type SocketMessage struct {
	From        string
	To          string
	MessageID   string // for reply
	MessageType int32
	State       Common.ErrorCodes //
	DataJson    string            // nullable
}

type RedisMessageReplyCheck struct {
	MessageID string
	Replied   bool
}

func Init() {
	UpdatePicbotTokenMap()

	// start ws
	http.HandleFunc("/ws", handleWebSocket)
	port := "8100"
	fmt.Printf("WebSocket server is listening on :%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func authenticateToken(r *http.Request) (bool, string) {
	// 获取 Authorization 头部的值
	authHeader := r.Header.Get("Authorization")
	utils.PrintObj(authHeader, "authHeader")
	// dao.CreateLog("authHeader", authHeader)

	// check token
	picbotName, exists := PicbotTokenMap[authHeader]
	utils.PrintObj(exists, "PicbotTokenMap exist")

	if !exists {
		utils.PrintObj(authHeader, "token not found")
		return false, ""
	}

	utils.PrintObj(picbotName, "connected")

	// check if aleady connected
	_, exists = WebSocketConns[picbotName]
	if exists {
		// dao.CreateLog("ws token already using", picbotName)
		utils.PrintObj(picbotName, "ws token already using")
		return false, ""
	}

	setPicbotAlive(picbotName, true)
	dao.CreateLog("setPicbotAliveState:"+picbotName, "true")

	return true, picbotName
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 在這裡打印客戶端的IP地址
	clientIP := r.RemoteAddr
	fmt.Printf("Client IP: %s\n", clientIP)

	// 验证令牌
	pass, reqPicbotName := authenticateToken(r)

	if !pass {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	WebSocketConns[reqPicbotName] = conn // 将连接添加到连接列表

	for {
		_, messageRaw, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		utils.PrintObj(string(messageRaw), "Received messgae")

		// 解析消息类型
		message, err := utils.ParseJsonWithType[SocketMessage](string(messageRaw))
		if err != nil {
			dao.CreateLog("Error parsing message type", err.Error())
			break
		}

		commandType, err := service.ConvertInt32ToPicbotCommandType(message.MessageType)
		if err != nil {
			utils.PrintObj(err.Error(), "err")
		}

		// 根据消息类型分发到不同的处理函数
		switch commandType {
		case PicbotEnum.PicbotCommandType_REPLY_FROM_PICBOT:
			handleReply(message)
		case PicbotEnum.PicbotCommandType_PICBOT_STATE:
			HandlePicbotState(message)
		case PicbotEnum.PicbotCommandType_ADD_ORDER:
			dao.CreateLog("Received message PicbotCommandType_ADD_ORDER", string(messageRaw))

			err := HandleAddOrderRecord(message)
			if err != nil {
				dao.CreateLog("HandleAddOrderRecord err", err.Error())
			}

		case PicbotEnum.PicbotCommandType_GET_THEME_GROUP_INFO:
			// dao.CreateLog("Received message PicbotCommandType_GET_THEME_GROUP_INFO", string(messageRaw))

			HandleGetThemeGroup(message)
		default:
			fmt.Println("Unknown message type:", message)
			return
		}

		// send confirm receive message
		receiveConfirm := SocketMessage{
			From:        BackendSenderNameDefault,
			To:          message.From,
			MessageID:   message.MessageID,
			MessageType: int32(PicbotEnum.PicbotCommandType_RECEIVE_CONFIRM),
		}

		code := SendMessage(receiveConfirm, false)

		if code != Common.ErrorCodes_SUCCESS {
			dao.CreateLog("SendMessage", code.String())
			return
		}
	}

	// 離線
	setPicbotAlive(reqPicbotName, false)
	delete(WebSocketConns, reqPicbotName)
	dao.CreateLog("setPicbotAliveState:"+reqPicbotName, "false")
}

func HandleGetThemeGroup(message SocketMessage) {
	// get picbot data
	code, _ := service.GetPicbotInfoForPicbot(message.From)
	if code != Common.ErrorCodes_SUCCESS {
		SendBackendReply(message.From, message.MessageID, ginResult.CommonResult{Code: code})
	}

	// mapping theme group data
	// res := service.GetThemeGroupForPicbot(picbot.ThemeGroupId)
	// SendBackendReply(message.From, message.MessageID, res)
}
