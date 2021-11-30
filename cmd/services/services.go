package services

import (
	"math/rand"
	"time"
)

// GenerateRandString – генератор строки из случайных символов(букв, чисел и _), len – длина строки
func GenerateRandString(len int) string {
	rand.Seed(time.Now().UnixNano())
	bStr := make([]byte, len)
	for i := range bStr {
		bStr[i] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"[rand.Intn(63)]
	}
	return string(bStr)
}
