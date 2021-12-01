package services

import (
	"encoding/json"
	"log"
	"main/cmd/models"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type IServices interface {
	Check() (bool, error)
	Create() error
}

var SERVER = os.Getenv("SERVER_NAME") + os.Getenv("SERVER_PORT")

type Response struct {
	Code        int    `json:"-"`
	Url         string `json:"url"`
	Error       bool   `json:"error"`
	Description string `json:"description"`
}

type Link struct {
	db  *models.Database
	r   *http.Request
	rsp *Response
}

type URL struct {
	Url string `json:"url"`
}

func (l *Link) LinkExecute() (*Response, error) {
	var resp = new(Response)
	if l.r.Method != http.MethodGet && l.r.Method != http.MethodPost {
		resp.Code = http.StatusMethodNotAllowed
		resp.Error = true
		resp.Description = "Method Not Allowed"
		return resp, nil
	}
	switch l.r.Method {
	case http.MethodPost:
		l.Post()
	case http.MethodGet:

	}
	return l.rsp, nil
}

func (l *Link) Get() {

}

func (l *Link) Post() {
	var url URL
	if err := json.NewDecoder(l.r.Body).Decode(&url); err != nil {
		l.rsp.Error = true
		if err.Error() == "EOF" {
			l.rsp.Code = http.StatusBadRequest
			l.rsp.Description = "Empty Body"
		} else {
			l.rsp.Code = http.StatusInternalServerError
			l.rsp.Description = "Internal Server Error"
		}
		return
	}
	if url.Url == "" {
		l.rsp.Code = http.StatusBadRequest
		l.rsp.Error = true
		l.rsp.Description = "Bad Request"
		return
	}
	// tmp
	l.rsp.Code = http.StatusCreated
	l.rsp.Url = url.Url
	l.rsp.Error = false
	l.rsp.Description = "Success"
	log.Println(l.r.Host)
}

func (l *Link) Check() (*Response, error) {
	return nil, nil
}

func NewLink(db *models.Database, r *http.Request) *Link {
	return &Link{db: db, r: r, rsp: &Response{}}
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
