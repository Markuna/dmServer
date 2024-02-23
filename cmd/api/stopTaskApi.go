package api

import (
	"bytes"
	"douyinApi/cmd/api/common"
	"douyinApi/config"
	"encoding/json"
	"fmt"
	"net/http"
)

func StopTask(req *common.TaskReq) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	url := config.Get().Douyin.Url + "/api/live_data/task/stop"
	jsonStr, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("access-token", GetAccessToken())

	rspBody := common.Do(httpReq)

	return rspBody.ErrNo == common.ErrCodeSuccess
}
