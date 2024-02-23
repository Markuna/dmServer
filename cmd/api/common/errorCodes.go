package common

import "fmt"

// 定义错误码
const (
	ErrCodeSuccess                      = 0       // 成功
	ErrCodeInvalidAppId                 = 40015   // appid 错误
	ErrCodeInvalidSecret                = 40017   // secret 错误
	ErrCodeInvalidGrantType             = 40020   // grant_type 不是 client_credential
	ErrCodeRequiredParamentersAreAbsent = 40023   // 缺少必要的请求参数
	ErrCodeInvalidAccessToken           = 40022   // 无效的 access_token
	ErrCodeServiceUnavaliable           = 10001   // 服务内部错误，调用下游服务失败
	ErrCodePushTaskCanNotStart          = 5003019 // 推送任务不符合开启的限制条件
	ErrCodeInternalServerError          = -1      // 服务器内部错误
)

// 定义错误信息
var ErrMsg = map[int]string{
	ErrCodeSuccess:                      "成功",
	ErrCodeInvalidAppId:                 "参数无效",
	ErrCodeInvalidSecret:                "权限不足",
	ErrCodeInvalidGrantType:             "grant_type 不是 client_credential",
	ErrCodeRequiredParamentersAreAbsent: "缺少必要的请求参数",
	ErrCodeInvalidAccessToken:           "无效的 access_token",
	ErrCodeServiceUnavaliable:           "服务内部错误，调用下游服务失败",
	ErrCodePushTaskCanNotStart:          "推送任务不符合开启的限制条件",
	ErrCodeInternalServerError:          "服务器内部错误",
}

// 定义一个错误结构体
type Error struct {
	Code    int    // 错误码
	Message string // 错误信息
}

// 实现error接口
func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
