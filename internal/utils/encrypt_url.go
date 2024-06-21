package utils

import (
	"fmt"
	"strconv"
	"time"
)

const (
	encryptKey = "DT06W167NBuTPg3GNvxjvfTo0jJ8lOOjk8pXd9WwR2w"
)

// EncryptURL https://example.com/xxx/xxx?sign=xxxxx&t=xxxx
func EncryptURL(uPath string, timeout int64) (string, string) {
	var (
		now      = time.Now().Unix()
		nowValue = fmt.Sprintf("%d", now/timeout)
	)
	str := encryptKey + uPath + nowValue
	return GetMD5Hash(str), Time16(now)
}

// DecryptURL 是否合法
func DecryptURL(uPath, nowParam, sign string, timeout int64) bool {
	var (
		now      = time.Now().Unix()
		nowValue = fmt.Sprintf("%d", now/timeout)
	)

	if now-time10(nowParam) > timeout {
		return false
	}

	return sign == GetMD5Hash(encryptKey+uPath+nowValue)
}

func Time16(now int64) string {
	return fmt.Sprintf("%x", now)
}

func time10(value string) int64 {
	now, _ := strconv.ParseInt(value, 16, 64)
	return now
}
