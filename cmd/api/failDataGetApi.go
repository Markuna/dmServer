package api

import (
	"douyinApi/cmd/api/common"
	"douyinApi/config"
	"encoding/json"
	"fmt"
	"net/http"
	urlApi "net/url"
	"strconv"
)

type FailDataReq struct {
	RoomId   string `json:"roomid"`
	AppId    string `json:"appid"`
	MsgType  string `json:"msg_type"`
	PageNum  int64  `json:"page_num"`
	PageSize int64  `json:"page_size"`
}
type DataListResp struct {
	PageNum    int64         `json:"page_num"`
	TotalCount int64         `json:"total_count"`
	DataList   []DataListStr `json:"data_list"`
}

type DataListData struct {
	PageNum    int64      `json:"page_num"`
	TotalCount int64      `json:"total_count"`
	DataList   []DataList `json:"data_list"`
}
type DataList struct {
	RoomId  string            `json:"room_id"`
	MsgType string            `json:"msg_type"`
	Payload []LiveGiftPayload `json:"payload"`
}

type DataListStr struct {
	RoomId  string `json:"room_id"`
	MsgType string `json:"msg_type"`
	Payload string `json:"payload"`
}

func GetFailData(req FailDataReq) *DataListData {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	url := config.Get().Douyin.Url + "/api/live_data/task/fail_data/get"
	values := urlApi.Values{}
	values.Set("roomid", req.RoomId)
	values.Set("appid", req.AppId)
	values.Set("msg_type", req.MsgType)
	values.Set("page_num", strconv.Itoa(int(req.PageNum)))
	values.Set("page_size", strconv.Itoa(int(req.PageSize)))
	httpReq, _ := http.NewRequest("GET", url+"?"+values.Encode(), nil)

	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("access-token", GetAccessToken())

	rspBody := common.Do(httpReq)

	var resp DataListResp
	s, _ := json.Marshal(rspBody.Data)
	json.Unmarshal(s, &resp)

	result := DataListData{
		PageNum:    resp.PageNum,
		TotalCount: resp.TotalCount,
	}
	result.DataList = make([]DataList, 0, len(resp.DataList))
	for _, v := range resp.DataList {
		var payload []LiveGiftPayload
		json.Unmarshal([]byte(v.Payload), &payload)
		result.DataList = append(result.DataList, DataList{
			RoomId:  v.RoomId,
			MsgType: v.MsgType,
			Payload: payload,
		})
	}
	return &result
}

type LiveGiftPayload struct {
	MsgId     string `json:"msg_id"`      // string类型id
	SecOpenId string `json:"sec_openid"`  // 评论用户的加密openid, 当前其实没有加密
	SecGiftId string `json:"sec_gift_id"` // 加密的礼物id
	GiftNum   int32  `json:"gift_num"`    // 送出的礼物数量
	GiftValue int32  `json:"gift_value"`  // 礼物总价值，单位分
	AvatarUrl string `json:"avatar_url"`  // 用户头像
	Nickname  string `json:"nickname"`    // 用户昵称(不加密)
	Timestamp int64  `json:"timestamp"`   // 时间戳
}
