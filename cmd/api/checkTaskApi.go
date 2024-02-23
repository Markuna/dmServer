package api

import (
	"bytes"
	"douyinApi/cmd/api/common"
	"douyinApi/config"
	"encoding/json"
	"fmt"
	"net/http"
)

type Status struct {
	Status int32 `json:"status"` // 1 任务不存在 2 任务未启动 3 任务运行中
}

func CheckTask(req *common.TaskReq) *Status {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	url := config.Get().Douyin.Url + "/api/live_data/task/get"
	jsonStr, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("access-token", GetAccessToken())

	rspBody := common.Do(httpReq)

	d := &Status{
		Status: rspBody.Data["status"].(int32),
	}
	return d
}
