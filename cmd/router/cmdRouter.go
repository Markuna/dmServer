package router

import (
	"douyinApi/cmd/api"
	"douyinApi/cmd/api/common"
	"douyinApi/cmd/service"
	"douyinApi/config"
	"douyinApi/consumer"
	"douyinApi/dao"
	"douyinApi/redis"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// Ping test
	r.GET("/pingCmd", func(c *gin.Context) {
		// go topGift("")
		c.String(http.StatusOK, "pongCmd")
	})

	// websocket连接数
	r.GET("/wsInfo", func(c *gin.Context) {
		size := consumer.WebsocketManager.LenGroup()
		roomIds := consumer.WebsocketManager.GroupNames()
		data := make(map[string]interface{})
		data["size"] = size
		data["roomIds"] = roomIds
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "",
			"data": data,
		})
	})

	// 开启链接
	r.POST("/login", func(ctx *gin.Context) {

		// 通过token 拿 roomId
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			panic(err)
		}
		var dmap map[string]interface{}
		err = json.Unmarshal(bodyBytes, &dmap)
		if err != nil {
			panic(err)
		}
		token := dmap["token"].(string)
		resp := api.RoomIdApi(&api.RoomIdTokenReq{
			Token: token,
		})
		if resp.ErrMsg != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  resp.ErrMsg,
				"data": nil,
			})
			return
		} else {
			//resp.Data.Info.RoomId -> dmap["token"].(string) set into redis
			key := strconv.FormatInt(resp.Data.Info.RoomId, 10)
			redis.RedisDb.Set(key, token, 0)
			startTasks(key)
			go topGift(key)
			ctx.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "",
				"data": gin.H{
					"roomId": resp.Data.Info.RoomId,
				},
			})
		}
	})

	// 获取世界排行榜
	r.GET("/wordRank", func(ctx *gin.Context) {
		pageNoStr := ctx.Query("pageNo")
		pageNo, _ := strconv.ParseInt(pageNoStr, 10, 32)
		pageSizeStr := ctx.Query("pageSize")
		pageSize, _ := strconv.ParseInt(pageSizeStr, 10, 32)
		resp := service.GetWordRankData(&service.PageParam{
			PageNo:   pageNo,
			PageSize: pageSize,
		})
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "",
			"data": resp,
		})
	})

	// 接收当局比分
	r.POST("/pushScopes", func(ctx *gin.Context) {
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			panic(err)
		}
		var darr common.ArrayData
		err = json.Unmarshal(bodyBytes, &darr)
		if err != nil {
			panic(err)
		}
		service.UpdateWordRank(darr.Data)
		resp := service.Get10UserWordRankInfo(darr.Data)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "",
			"data": resp,
		})
	})

	// 开始日志
	r.POST("/log/startInfo", func(ctx *gin.Context) {
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			panic(err)
		}
		var d service.StartInfo
		err = json.Unmarshal(bodyBytes, &d)
		if err != nil {
			panic(err)
		}
		d.Save()
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "",
			"data": nil,
		})
	})

	// 对局日志
	r.POST("/log/gameInfo", func(ctx *gin.Context) {
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			panic(err)
		}
		var d service.GameInfo
		err = json.Unmarshal(bodyBytes, &d)
		if err != nil {
			panic(err)
		}
		d.Save()
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "",
			"data": nil,
		})
	})

	// 玩家日志
	r.POST("/log/playerInfo", func(ctx *gin.Context) {
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			panic(err)
		}
		var d service.PlayerInfo
		err = json.Unmarshal(bodyBytes, &d)
		if err != nil {
			panic(err)
		}
		d.Save()
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "",
			"data": nil,
		})
	})

	// 白名单查询
	r.GET("/whiteList/check", func(ctx *gin.Context) {
		uid := ctx.Query("uid")
		resp := service.CheckUidIsInWhiteList(uid)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "",
			"data": resp,
		})
	})

	// 通用字典查询
	r.GET("/dict/query", func(ctx *gin.Context) {
		uid := ctx.Query("key")
		resp := service.QueryDictByKey(uid)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "",
			"data": resp,
		})
	})
}

func startTasks(roomId string) {
	resp := api.StartTask(&common.TaskReq{
		Roomid:  roomId,
		Appid:   config.Get().Douyin.AppId,
		MsgType: api.LiveComment,
	})
	log.Println(resp.ErrMsg)
	log.Println(resp.Data.TaskId)
	resp = api.StartTask(&common.TaskReq{
		Roomid:  roomId,
		Appid:   config.Get().Douyin.AppId,
		MsgType: api.LiveGift,
	})
	log.Println(resp.ErrMsg)
	log.Println(resp.Data.TaskId)
	resp = api.StartTask(&common.TaskReq{
		Roomid:  roomId,
		Appid:   config.Get().Douyin.AppId,
		MsgType: api.LiveLike,
	})
	log.Println(resp.ErrMsg)
	log.Println(resp.Data.TaskId)
}

func topGift(roomId string) {
	giftIds := []string{}
	giftIdsDict := dao.QueryDictByKey("topGift")
	for _, v := range *giftIdsDict {
		giftIds = append(giftIds, v.DValue)
	}
	resp := api.TopGift(&api.TopGiftReq{
		RoomId:        roomId,
		AppId:         config.Get().Douyin.AppId,
		SecGiftIdList: giftIds,
	})
	if resp.ErrMsg != "" {
		log.Println(resp.ErrMsg)
	} else {
		log.Println(resp.Data)
	}
}
