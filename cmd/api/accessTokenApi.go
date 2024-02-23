package api

import (
	"bytes"
	"douyinApi/cmd/api/common"
	"douyinApi/config"
	"encoding/json"
	"fmt"
	"net/http"
)

type AccessTokenReq struct {
	AppId      string `json:"appId"`      // 小程序id
	Secret     string `json:"secret"`     // 小程序的App Secret
	Grant_type string `json:"grant_type"` // 获取 access_token 时为 client_credential
}

type AccessTokenRespData struct {
	AccessToken string  `json:"access_token"` // access_token
	ExpiresIn   float64 `json:"expires_in"`   // 过期时间 秒
}

// Tip: token 是小程序级别 token，不要为每个用户单独分配一个 token，会导致 token 校验失败。建议每小时更新一次即可。
func AccessTokenApi(req *AccessTokenReq) *AccessTokenRespData {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	url := config.Get().Douyin.ToutiaoUrl + "/api/apps/v2/token"
	jsonStr, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")

	rspBody := common.Do(httpReq)

	d := &AccessTokenRespData{
		AccessToken: rspBody.Data["access_token"].(string),
		ExpiresIn:   rspBody.Data["expires_in"].(float64),
	}
	return d
}
