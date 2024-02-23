package common

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func Do(httpReq *http.Request) *CommonResponse {

	out := DoNative(httpReq)
	rspBody := new(CommonResponse)
	err := json.Unmarshal(out, rspBody)
	if err != nil {
		panic(err)
	}
	if rspBody.ErrNo != ErrCodeSuccess {
		panic(errors.New(rspBody.ErrTips))
	}
	return rspBody
}

func DoNative(httpReq *http.Request) []byte {

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic(errors.New("http Failed, code is not 200 url:" + httpReq.RequestURI + " code:" + resp.Status))
	}
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if out == nil {
		panic(errors.New("http Failed, body is empty, url:" + httpReq.RequestURI))
	}
	return out
}
