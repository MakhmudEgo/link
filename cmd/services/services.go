package services

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"log"
	"main/cmd/models"
	"main/cmd/utils"
	"math/rand"
	"net/http"
	"time"
)

type IServices interface {
	Select(string) (*models.Link, error)
	Insert(string) (*models.Link, error)
	Create(string) (*models.Link, error)
}

type Response struct {
	Code        int    `json:"-"`
	Url         string `json:"url"`
	Error       bool   `json:"error"`
	Description string `json:"description"`
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
		l.Get()
	}
	return l.rsp, nil
}

func (l *Link) Get() {
	l.rsp.Code = http.StatusTeapot
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
	res, err := l.Select(url.Url)
	if err != pgx.ErrNoRows && err != redis.Nil && err != nil {
		log.Println(err)
		l.rsp.Code = http.StatusInternalServerError
		l.rsp.Description = "Internal Server Error"
	} else if err != nil {
		l.rsp.Code = http.StatusCreated
		res, _ = l.Create(url.Url) //err
		l.rsp.Description = "Created"
	} else {
		l.rsp.Code = http.StatusOK
		l.rsp.Description = "Success"
	}
	l.rsp.Url = res.Short
	l.rsp.Error = false

	log.Println(l.r.Host)
}

func (l *Link) Create(url string) (*models.Link, error) {
	link := &models.Link{}
	link.Base = url
	var err error
	short := l.SERVER_URL + GenerateRandString(10)
	switch l.db.TypeDB {
	case utils.PostgresDB:
		for err = l.db.Pst.QueryRow(
			context.Background(),
			"insert into link(base, short) values ($1, $2)",
			url, short).Scan(); err != pgx.ErrNoRows; {
			short = l.SERVER_URL + GenerateRandString(10)
		}

		link.Short = short

	case utils.RedisDB:
		//link.Short, err = l.db.Rdb.Get(l.db.Ctx, url).Result()
	}
	return link, err
}

func (l *Link) Select(url string) (*models.Link, error) {
	link := &models.Link{}
	link.Base = url
	var err error
	switch l.db.TypeDB {
	case utils.PostgresDB:
		err = l.db.Pst.QueryRow(
			context.Background(),
			"select short from link where base=$1",
			url).Scan(&link.Short)
	case utils.RedisDB:
		link.Short, err = l.db.Rdb.Get(l.db.Ctx, url).Result()
	}
	return link, err
}

func NewLink(db *models.Database, r *http.Request, SERVER_URL string) *Link {
	return &Link{db: db, r: r, rsp: &Response{}, SERVER_URL: SERVER_URL}
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
