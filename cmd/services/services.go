package services

import (
	"main/cmd/models"
	"math/rand"
	"net/http"
	"time"
)

type IServices interface {
	Select(string) (*models.Link, error)
	//Create(string) (*models.Link, error)
}

type Response struct {
	Code        int    `json:"-"`
	Url         string `json:"url"`
	Error       bool   `json:"error"`
	Description string `json:"description"`
}

type Home struct {
	db         *models.Database
	r          *http.Request
	SERVER_URL string
}

type Link struct {
	db         *models.Database
	r          *http.Request
	rsp        *Response
	SERVER_URL string
}


type URL struct {
	Url string `json:"url"`
}

// GenerateRandString – генератор строки из случайных символов(букв, чисел и _), len – длина строки
func GenerateRandString(len int) string {
	rand.Seed(time.Now().UnixNano())
	bStr := make([]byte, len)
	for i := range bStr {
		bStr[i] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"[rand.Intn(63)]
	}
	return string(bStr)
}
