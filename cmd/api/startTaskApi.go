package api

import (
	"bytes"
	"douyinApi/cmd/api/common"
	"douyinApi/config"
	"encoding/json"
	"log"
	"net/http"
)

const (
	LiveComment = "live_comment" // 评论
	LiveGift    = "live_gift"    // 礼物
	LiveLike    = "live_like"    // 点赞
)

type RespTask struct {
	Data   TaskId `json:"data"`
	ErrMsg string `json:"err_msg"`
	ErrNo  int    `json:"err_no"`
	Logid  string `json:"logid"`
}
type TaskId struct {
	TaskId string `json:"task_id"` // 任务id，每次启动任务都会有一个
}

func StartTask(req *common.TaskReq) *RespTask {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	url := config.Get().Douyin.Url + "/api/live_data/task/start"
	jsonStr, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("access-token", GetAccessToken())

	rspBody := common.DoNative(httpReq)
	var d RespTask
	err := json.Unmarshal(rspBody, &d)
	if err != nil {
		panic(err)
	}
	return &d
}
