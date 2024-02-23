package main

import (
	callback "douyinApi/callback/router"
	cmd "douyinApi/cmd/router"
	"douyinApi/config"
	"douyinApi/consumer"
	"log"

	_ "douyinApi/cmd/api"
	_ "douyinApi/dao"
	_ "douyinApi/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	httpServer()
}

func httpServer() {

	// 路由初始化
	gin.SetMode(gin.ReleaseMode)
	rx := gin.Default()
	callback.SetupRouter(rx)
	cmd.SetupRouter(rx)

	// websocket
	go consumer.WebsocketManager.Start()
	go consumer.WebsocketManager.SendService()

	wsGroup := rx.Group("/ws")
	wsGroup.GET("/:channel", consumer.WebsocketManager.WsClient)

	// Listen and Server in 0.0.0.0:8080
	rx.Run(":" + config.Get().Server.Http.Port)
}
