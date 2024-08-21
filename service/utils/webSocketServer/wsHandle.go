package webSocketServer

import (
	"errors"
	"fmt"
	"service/dao"
	"service/model"
	"service/protos/Common"
	"service/protos/PicbotEnum"
	"service/service"
	"service/utils"
	"service/utils/ginResult"
	"service/utils/redisClient"
)

type AddOrderRequest struct {
	Order       model.Orders
	OrderDetail model.OrderDetails
}

type ImageToImageResult struct {
	PicbotProductId string
	Url             string
}

func HandleAddOrderRecord(message SocketMessage) error {
	orderRequest, err := utils.ParseJsonWithType[AddOrderRequest](message.DataJson)

	if err != nil {
		return err
	}

	dao.CreateLog("HandleAddOrderRecord", utils.ToJson(orderRequest))

	// check price
	validPrice := (orderRequest.Order.ItemPrice * orderRequest.Order.AmountPrintOut) == orderRequest.Order.TotalPrice

	// check carrier
	validCarrier := true
	if orderRequest.Order.CarruerNumber != "" {
		validCarrier = utils.IsValidCarrierBarcode(orderRequest.Order.CarruerNumber)
	}

	dao.CreateLog("HandleAddOrderRecord check: validPrice, validCarrier", utils.ToJson([]bool{validPrice, validCarrier}))

	// check
	if !validPrice ||
		orderRequest.Order.TotalPrice == 0 ||
		!validCarrier {
		return errors.New("valid fail")
	}

	// add order into DB
	if orderRequest.Order.PicbotName != "" &&
		orderRequest.Order.PaymentType != PicbotEnum.PicbotPaymentType_PicbotPaymentType_NONE &&
		orderRequest.Order.State != PicbotEnum.PicbotOrderState_PicbotOrderState_NONE {

		// add order
		orderRequest.Order.InvoiceState = PicbotEnum.PicbotOrderInvoiceState_INVOICE_PEDDING
		res, order := dao.Create(orderRequest.Order)

		if res.Error != nil {
			return res.Error
		}
		utils.PrintObj(order, "create order data")

		// mapping order id and create detail
		orderRequest.OrderDetail.OrderId = order.Id

		res, orderDetail := dao.Create(orderRequest.OrderDetail)

		if res.Error != nil {
			return res.Error
		}

		utils.PrintObj(orderDetail, "create order detail data")
		dao.CreateLog("create order", "success")
		return nil
	}

	return errors.New("PicbotName or PaymentType or State is invalid")
}

func SendBackendReply(to, messageId string, replyData ginResult.CommonResult) {
	receiveConfirm := SocketMessage{
		From:        BackendSenderNameDefault,
		To:          to,
		MessageID:   messageId,
		MessageType: int32(PicbotEnum.PicbotCommandType_REPLY_FROM_BACKEND),
		DataJson:    utils.ToJson(replyData),
	}

	code := SendMessage(receiveConfirm, false)
	if code != Common.ErrorCodes_SUCCESS {
		dao.CreateLog("SendMessage", code.String())
		return
	}
}

type StateChangeRequest struct {
	Field string
	Value string
}

func HandlePicbotState(message SocketMessage) {

	data, err := utils.ParseJsonWithType[StateChangeRequest](message.DataJson)

	if err != nil {
		dao.CreateLog("handlePicbotState:"+err.Error(), utils.ToJson(message))
		return
	}

	picbotName := message.From

	// trying to get state from redis if cant than create new one
	state := service.PicbotState{}
	redisKey := redisClient.RedisKeyPicbotState + picbotName
	res := redisClient.GetRedis(redisKey)

	if res != "" {
		data, err := utils.ParseJsonWithType[service.PicbotState](res)

		if err != nil {
			dao.CreateLog("handlePicbotState", err.Error())
		} else {
			state = data
		}
	}

	state.PicbotName = picbotName
	update := false

	// valid and add new
	switch data.Field {
	case "PicbotBootTime":
		if utils.IsValidDate(data.Value) {
			state.PicbotBootTime = data.Value
			update = true
		}

	case "Temperature":
		parseResult := utils.ToInt(data.Value)

		if parseResult != -1 {
			state.Temperature = parseResult
			update = true
		}

	case "RemainingPages":
		parseResult := utils.ToInt(data.Value)

		if parseResult != -1 {
			state.RemainingPages = parseResult
			update = true
		}

	case "TakePictureCount":
		parseResult := utils.ToInt(data.Value)

		if parseResult != -1 {
			state.TakePictureCount = parseResult
			update = true
		}

	case "PrintCount":
		parseResult := utils.ToInt(data.Value)

		if parseResult != -1 {
			state.PrintCount = parseResult
			update = true
		}

	case "PeopleCount":
		parseResult := utils.ToInt(data.Value)

		if parseResult != -1 {
			state.PeopleCount = parseResult
			update = true
		}
	case "NetworkSpeedUpload":
		parseResult := utils.ToInt(data.Value)

		if parseResult != -1 {
			state.NetworkSpeedUpload = parseResult
			update = true
		}
	case "NetworkSpeedDownload":
		parseResult := utils.ToInt(data.Value)

		if parseResult != -1 {
			state.NetworkSpeedDownload = parseResult
			update = true
		}
	default:
		dao.CreateLog("unknown type", data.Field)
		return
	}

	if update {
		redisClient.SetRedisForever(redisKey, utils.ToJson(state))
	} else {
		dao.CreateLog("wrong type, not update:"+data.Field, data.Value)
	}
}

func handleReply(message SocketMessage) {
	// dao.CreateLogJson("handleReply", message)

	// valid
	if !utils.ValidId(message.MessageID) {
		utils.PrintObj(message.MessageID, "mqttMessage.MessageID error")
		return
	}

	messageCheck := RedisMessageReplyCheck{
		MessageID: message.MessageID,
		Replied:   true,
	}

	// update redis
	key := redisClient.RedisKeySocketWaitReplyMessages + message.MessageID
	res := redisClient.GetRedis(key)
	if res != "" {
		redisClient.SetRedis(key, utils.ToJson(messageCheck), 1)
		// dao.CreateLog("handleReplyMessage ok", message.MessageID)
	} else {
		fmt.Println("handleReply => cant find redis:" + key)
	}
}
