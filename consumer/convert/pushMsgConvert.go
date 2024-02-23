package convert

import (
	"douyinApi/consumer/proto"
	"encoding/json"
	"log"
)

func ConvertPushMsg(msgType string, body []byte) *proto.PushMsg {

	return &proto.PushMsg{
		PushType: msgType,
		Payload:  convert2Payload(body),
	}
}

func convert2Payload(body []byte) *proto.Payload {
	var payload proto.Payload
	err := json.Unmarshal([]byte(body), &payload)
	if err != nil {
		log.Println(string(body))
		log.Println(err)
	}
	return &payload
}
