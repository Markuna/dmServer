package utils

import (
	"crypto/md5"
	"encoding/base64"
	"sort"
	"strings"
)

/*
*
比如：

	header := map[string]string{
	             "x-nonce-str": "123456",
	             "x-timestamp": "456789",
	             "x-roomid":    "268",
	             "x-msg-type":  "live_gift",
	     }
	bodyStr := "abc123你好"
	secret := "123abc"

rawData为：x-msg-type=live_gift&x-nonce-str=123456&x-roomid=268&x-timestamp=456789abc123你好123abc
signature为：PDcKhdlsrKEJif6uMKD2dw==
*/
func Signature(header map[string]string, bodyStr, secret string) string {
	keyList := make([]string, 0, 4)
	for key := range header {
		keyList = append(keyList, key)
	}
	sort.Slice(keyList, func(i, j int) bool {
		return keyList[i] < keyList[j]
	})
	kvList := make([]string, 0, 4)
	for _, key := range keyList {
		kvList = append(kvList, key+"="+header[key])
	}
	urlParams := strings.Join(kvList, "&")
	rawData := urlParams + bodyStr + secret
	md5Result := md5.Sum([]byte(rawData))
	return base64.StdEncoding.EncodeToString(md5Result[:])
}
