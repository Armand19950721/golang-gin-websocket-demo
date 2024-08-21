package webSocketServer

import (
	"errors"
	"fmt"
	"service/dao"
	"service/model"
	"service/protos/Common"
	"service/service"
	"service/utils"
	"service/utils/redisClient"
	"time"

	"github.com/gorilla/websocket"
)

// not use for now
func BroadcastMessage(message SocketMessage, waitReply bool) Common.ErrorCodes {
	for _, conn := range WebSocketConns {
		err := conn.WriteMessage(websocket.TextMessage, []byte(utils.ToJson(message)))
		if err != nil {
			dao.CreateLog("Error broadcast message", err.Error())
			return Common.ErrorCodes_INTERNAL_ERROR
		}
	}

	// 等回應
	if waitReply {
		return WaitReply(message.MessageID)
	}

	return Common.ErrorCodes_SUCCESS
}

func SendMessage(message SocketMessage, waitReply bool) Common.ErrorCodes {
	// check if aleady connected
	conn, exists := WebSocketConns[message.To]
	if !exists {
		dao.CreateLog("SendMessage", "receiver picbot is not in connection")
		return Common.ErrorCodes_RECEIVER_NOT_EXIST_IN_CONNECTION_POOL
	}

	// send message
	err := conn.WriteMessage(websocket.TextMessage, []byte(utils.ToJson(message)))
	if err != nil {
		dao.CreateLog("Error send message", err.Error())
		return Common.ErrorCodes_INTERNAL_ERROR
	}

	// 等回應
	if waitReply {
		return WaitReply(message.MessageID)
	}

	return Common.ErrorCodes_SUCCESS
}

func WaitReply(newMessageID string) Common.ErrorCodes {
	if newMessageID == "" {
		return Common.ErrorCodes_INTERNAL_ERROR
	}

	// 存redis messageId
	redisData := utils.ToJson(RedisMessageReplyCheck{
		MessageID: newMessageID,
		Replied:   false,
	})

	err := redisClient.SetRedis(redisClient.RedisKeySocketWaitReplyMessages+newMessageID, redisData, 1)
	if err != nil {
		fmt.Println(err.Error())
		return Common.ErrorCodes_INTERNAL_ERROR
	}

	// 等回應
	err = WaitingPicbotReply(newMessageID)
	utils.PrintObj(err, "WaitForPicbotReply err")

	if err != nil {
		fmt.Println(err.Error())
		return Common.ErrorCodes_SEND_COMMAND_TIMEOUT
	} else {
		return Common.ErrorCodes_SUCCESS
	}
}

func WaitingPicbotReply(messageId string) error {
	interval := 2 * time.Second
	count := 0
	countLimit := 3

	ticker := time.NewTicker(interval)
	for range ticker.C {
		// valid
		utils.PrintObj("retry check reply")

		if count >= countLimit {
			return errors.New("wait for " + utils.ToString(countLimit) + "second. and picbot didnt reply")
		}

		// check redis
		res := redisClient.GetRedis(redisClient.RedisKeySocketWaitReplyMessages + messageId)
		if res == "" {
			return errors.New("didnt create redis data at publish so is err")
		} else {
			// parse json
			data, err := utils.ParseJsonWithType[RedisMessageReplyCheck](res)
			if err != nil {
				return err
			}
			// check message state
			if data.Replied {
				utils.PrintObj(messageId, "received picbot reply")
				return nil
			}
		}

		count++
	}

	fmt.Printf("should not execute here")
	return nil
}

func setPicbotAlive(picbotName string, isAlive bool) {

	lastInfo := redisClient.GetRedis(redisClient.RedisKeyPicbotAlive + picbotName)

	connectTime := ""
	disconnecTime := ""

	if lastInfo != "" {
		// load old data
		res, err := utils.ParseJsonWithType[service.PicbotAlive](lastInfo)
		if err == nil {
			disconnecTime = res.DisconnectTime
			connectTime = res.ConnectTime
		}
	} else {
		// set picbot used
		res, _ := dao.Update(model.Picbot{Name: picbotName}, model.Picbot{Used: true})
		if res.Error != nil {
			dao.CreateLog("set picbot used:"+picbotName, res.Error.Error())
		}
	}

	if isAlive {
		connectTime = utils.Now()
	} else {
		disconnecTime = utils.Now()
	}

	state := service.PicbotAlive{
		PicbotName:     picbotName,
		IsAlive:        isAlive,
		ConnectTime:    connectTime,
		DisconnectTime: disconnecTime,
	}

	fmt.Println(utils.ToJson(state))

	redisClient.SetRedisForever(redisClient.RedisKeyPicbotAlive+picbotName, utils.ToJson(state))
}

func UpdatePicbotTokenMap() {
	// load picbot token map
	res, data, _ := dao.GetList(model.Picbot{})
	if res.Error != nil {
		dao.CreateLog("UpdatePicbotTokenMap", res.Error.Error())
	}

	for _, picbot := range data {
		PicbotTokenMap[picbot.SocketToken] = picbot.Name
	}

	for _, picbot := range data {
		PicbotTokenMap[picbot.SocketTokenOld] = picbot.Name
	}

	utils.PrintObj(PicbotTokenMap, "picbotTokenMap")
}
