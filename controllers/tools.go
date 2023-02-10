package controllers

import (
	"math/rand"
	"time"
)

//生成随机code
func GetHashCode(length int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, length)

	rand.Seed(time.Now().UnixNano())

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
