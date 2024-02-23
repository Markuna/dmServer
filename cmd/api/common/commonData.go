package common

import "douyinApi/cmd/service"

type CommonResponse struct {
	ErrNo   int32                  `json:"err_no"`
	ErrTips string                 `json:"err_tips"`
	LogId   string                 `json:"logid"`
	Data    map[string]interface{} `json:"data"`
}

type TaskReq struct {
	Roomid  string `json:"roomid"`   // 直播间id
	Appid   string `json:"appid"`    // 小玩法id
	MsgType string `json:"msg_type"` // 评论,礼物,点赞
}

type ArrayData struct {
	Data []service.WordRankData `json:"data"`
}
