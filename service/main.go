package main

import (
	"fmt"
	"service/server"
	"service/utils"
	"service/utils/postgresClient"
	"service/utils/redisClient"
	"service/utils/webSocketServer"
)

func main() {

	// load .env
	success := utils.InitEnv()
	if !success {
		return
	}

	// init DB
	success = postgresClient.InitDB()
	if !success {
		panic("InitDB: fail")
	} else {
		fmt.Println("InitDB: success")
	}

	// init redis
	success = redisClient.InitRedis()
	if !success {
		panic("InitRedis: fail")
	} else {
		fmt.Println("InitRedis: success")
	}

	// webSocket
	go webSocketServer.Init()
	fmt.Println("websocket: success")

	// start gin server
	server.InitGinServer()
}
