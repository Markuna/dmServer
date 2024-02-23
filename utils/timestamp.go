package utils

import (
	"strconv"
	"time"
)

func Timestamp() string {
	currentTime := time.Now().Unix()
	return strconv.FormatInt(currentTime, 10)
}
