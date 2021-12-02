package services

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"log"
	"main/cmd/models"
	"main/cmd/utils"
	"net/http"
	"strings"
)

func (l *Link) LinkExecute() (*Response, error) {
	if l.r.Method != http.MethodGet && l.r.Method != http.MethodPost {
		l.rsp.Code = http.StatusMethodNotAllowed
		l.rsp.Error = true
		l.rsp.Description = "Method Not Allowed"
		return l.rsp, nil
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
	var url = &URL{}
	if !l.IsCorrectRequest(url) {
		return
	}

	if !strings.HasPrefix(url.Url, l.SERVER_URL[:len(l.SERVER_URL)-1]) {
		l.rsp.Code = http.StatusBadRequest
		l.rsp.Error = true
		l.rsp.Description = "No Server URL Are Not Supported"
		return
	}

	res, err := l.Select(url.Url)
	if err != pgx.ErrNoRows && err != redis.Nil && err != nil {
		log.Println(err)
		l.rsp.Code = http.StatusInternalServerError
		l.rsp.Description = "Internal Server Error"
		l.rsp.Error = true
		return
	} else if err == pgx.ErrNoRows || err == redis.Nil {
		log.Println(err)
		l.rsp.Code = http.StatusNotFound
		l.rsp.Description = "Not Found"
		l.rsp.Error = true
		return
	}

	l.rsp.Code = http.StatusOK
	l.rsp.Url = res.Short
	l.rsp.Error = false
	l.rsp.Description = "Success"

	log.Println("services.go:", l.r.Host)
	log.Println("services.go, request URL:", url.Url)
}

func (l *Link) Post() {
	var url = &URL{}
	if !l.IsCorrectRequest(url) {
		return
	}

	if strings.HasPrefix(url.Url, l.SERVER_URL[:len(l.SERVER_URL)-1]) {
		l.rsp.Code = http.StatusBadRequest
		l.rsp.Error = true
		l.rsp.Description = "Server URL Are Not Supported"
		return
	}
	res, err := l.Select(url.Url)
	if err != pgx.ErrNoRows && err != redis.Nil && err != nil {
		log.Println(err)
		l.rsp.Code = http.StatusInternalServerError
		l.rsp.Description = "Internal Server Error"
		l.rsp.Error = true
		return
	} else if err != nil {
		l.rsp.Code = http.StatusCreated
		res, err = l.Create(url.Url)
		if err != nil && err != pgx.ErrNoRows {
			l.rsp.Code = http.StatusInternalServerError
			l.rsp.Description = "Internal Server Error"
			l.rsp.Error = true
			log.Println(err)
			return
		} else {
			l.rsp.Description = "Created"
		}
	} else {
		l.rsp.Code = http.StatusOK
		l.rsp.Description = "Success"
	}
	l.rsp.Url = res.Short
	l.rsp.Error = false

	log.Println("services.go:", l.r.Host)
	log.Println("services.go, request URL:", url.Url)
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
		err = l.db.Pst.QueryRow(
			context.Background(),
			"insert into link(base, short) values ($1, $2)",
			short, url).Scan()
	case utils.RedisDB:
		err = l.db.Rdb.Set(l.db.Ctx, url, short, 0).Err()
		if err != nil {
			return link, err
		}
		err = l.db.Rdb.Set(l.db.Ctx, short, url, 0).Err()
	}
	link.Short = short

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

func (l *Link) IsCorrectRequest(url *URL) bool {
	defer l.r.Body.Close()
	if err := json.NewDecoder(l.r.Body).Decode(&url); err != nil {
		l.rsp.Error = true
		if err.Error() == "EOF" {
			l.rsp.Code = http.StatusBadRequest
			l.rsp.Description = "Empty Body"
		} else {
			l.rsp.Code = http.StatusInternalServerError
			l.rsp.Description = "Internal Server Error"
		}
		return false
	}
	if url.Url == "" {
		l.rsp.Code = http.StatusBadRequest
		l.rsp.Error = true
		l.rsp.Description = "Bad Request"
		return false
	}
	return true
}
