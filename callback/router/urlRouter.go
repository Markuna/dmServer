package router

import (
	"douyinApi/callback/service"
	"douyinApi/config"
	"douyinApi/utils"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	// Ping test
	r.GET("/pingCallback", func(c *gin.Context) {
		c.String(http.StatusOK, "pongCallback")
	})

	// 参数形式 参数名 类型 说明
	// Header
	// x-nonce-str string 签名用的随机字符串
	// x-timestamp int64 发送消息的毫秒级时间戳
	// x-signature  string 请求签名，业务方接收后需要计算和校验签名，防伪造和篡改
	// x-roomid string 房间id
	// x-msg-type string 消息类型
	// 	live_comment:  直播间评论
	// 	live_gift: 直播间送礼
	// 	live_like: 直播间点赞
	// content-type  string 固定值：application/json
	// body ~   string 具体类型消息payload对象的json序列化字符串
	// 回调地址
	r.POST("/", func(c *gin.Context) {
		// header
		msgType := c.GetHeader("x-msg-type")
		roomId := c.GetHeader("x-roomid")
		header := make(map[string]string)
		header["x-msg-type"] = msgType
		header["x-roomid"] = roomId
		header["x-timestamp"] = c.GetHeader("x-timestamp")
		header["x-nonce-str"] = c.GetHeader("x-nonce-str")
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		bodyString := string(bodyBytes)

		sign := utils.Signature(header, bodyString, config.Get().Douyin.CallbackSecret)
		if sign != c.GetHeader("x-signature") && !config.Get().ForTest {
			log.Println("callback sign error")
			c.JSON(http.StatusOK, gin.H{"result": 0})
		} else {
			service.HandleMsg(msgType, roomId, bodyString)
			c.JSON(http.StatusOK, gin.H{"result": 1})
		}
	})
}
