package api

import (
	"bytes"
	"douyinApi/cmd/api/common"
	"douyinApi/config"
	"encoding/json"
	"log"
	"net/http"
)

type TopGiftReq struct {
	RoomId        string   `json:"room_id"`          // 直播间id
	AppId         string   `json:"app_id"`           // 小玩法id
	SecGiftIdList []string `json:"sec_gift_id_list"` // 顶置礼物id
}

type RespTopGift struct {
	Data    SecGiftIdList `json:"data"`
	ErrMsg  string        `json:"errmsg"`
	Errcode int           `json:"errcode"`
}
type SecGiftIdList struct {
	SecGiftIdList []string `json:"success_top_gift_id_list"`
}

func TopGift(req *TopGiftReq) *RespTopGift {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	url := config.Get().Douyin.Url + "/api/gift/top_gift"
	jsonStr, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("x-token", GetAccessToken())

	rspBody := common.DoNative(httpReq)
	var d RespTopGift
	err := json.Unmarshal(rspBody, &d)
	if err != nil {
		panic(err)
	}
	return &d
}
