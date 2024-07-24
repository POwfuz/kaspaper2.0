package kaspaperlib

import (
	"math/rand"
	"time"
)

// 生成16位 随机密码
func generatePassword() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < 16; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)

}
