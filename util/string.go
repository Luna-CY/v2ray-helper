package util

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

// Md5 计算MD5
func Md5(raw string) string {
	bytes := md5.Sum([]byte(raw))

	return fmt.Sprintf("%x", bytes)
}

// GenerateRandomString 随机生成字符串
func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	bytes := make([]rune, length)
	for i := range bytes {
		bytes[i] = runes[rand.Intn(len(runes))]
	}

	return string(bytes)
}
