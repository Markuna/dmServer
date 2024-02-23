package api

import (
	"bytes"
	"douyinApi/cmd/api/common"
	"douyinApi/config"
	"encoding/json"
	"fmt"
	"net/http"
)

type RoomIdTokenReq struct {
	Token string `json:"token"` // 直播伴侣给的token
}

type RoomIdRespData struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Data    struct {
		Info struct {
			RoomId int64 `json:"room_id"`
		} `json:"info"`
	} `json:"data"`
}

func RoomIdApi(req *RoomIdTokenReq) *RoomIdRespData {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	url := config.Get().Douyin.Url + "/api/webcastmate/info"
	jsonStr, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("X-Token", GetAccessToken())

	rspBody := common.DoNative(httpReq)

	var d RoomIdRespData
	err := json.Unmarshal(rspBody, &d)
	if err != nil {
		panic(err)
	}

	return &d
}
