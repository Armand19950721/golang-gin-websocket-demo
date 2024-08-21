package server

import (
	"log"
	"service/controller"
	"service/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitGinServer() {
	r := setRouter()
	err := r.Run(utils.GetEnv("GIN_SERVER_IP"))
	if err != nil {
		log.Fatalf("upload server run fail")
		return
	}
}

func setRouter() *gin.Engine {
	route := gin.Default()

	corsConf := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowHeaders: []string{"authorization", "Authorization", "Content-Type", "Upgrade", "Origin",
			"Connection", "Accept-Encoding", "Accept-Language", "Host", "Access-Control-Request-Method",
			"Access-Control-Request-Headers"},
	}
	// 設定靜態檔案路徑
	prefix := "/api/"

	route.Use(cors.New(corsConf))
	route.GET(prefix+"Echo/Echo", controller.Echo)

	return route
}
