package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"main/internal/pkg/utils"
	"main/pkg/controllers"
	"main/pkg/models"
	"main/pkg/services"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type Rsp struct {
	data    []byte
	rspCode int
	url     string
	rspUrl  string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}

func TestHome(t *testing.T) {
	SERVER_URL := os.Getenv("SERVER_URL")

	db := models.NewDatabase(utils.PostgresDB)
	if err := db.Connect(); err != nil {
		log.Fatalln(err)
	}
	svr := httptest.NewServer(controllers.NewHome(db, SERVER_URL))
	defer svr.Close()
	req, err := http.NewRequest("POST", svr.URL, nil)
	if err != nil {
		t.Error(err)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
}

func NewRspWithDataPost() []Rsp {
	var rsp = []Rsp{
		{data: []byte(`{"kek": "lol"}`), rspCode: http.StatusBadRequest},
		{data: nil, rspCode: http.StatusBadRequest},
		{data: []byte(`{"url": "https://ozon.ru"}`), rspCode: http.StatusCreated, url: "https://ozon.ru"},
		{data: []byte(`{"url": "https://ozon.ru"}`), rspCode: http.StatusOK, url: "https://ozon.ru"},
		{data: []byte(`{"url": "https://yandex.ru"}`), rspCode: http.StatusCreated, url: "https://yandex.ru"},
		{data: []byte(`{"url": "https://google.ru"}`), rspCode: http.StatusCreated, url: "https://google.ru"},
		{data: []byte(`{"url": "https://google.ru"}`), rspCode: http.StatusOK, url: "https://google.ru"},
		{data: []byte(`{"url": "https://yandex.ru"}`), rspCode: http.StatusOK, url: "https://yandex.ru"},
		{data: []byte(`{"url": "http://localhost:8080"}`), rspCode: http.StatusBadRequest, url: "http://localhost:8080"},
		{data: []byte(`{"url": ""}`), rspCode: http.StatusBadRequest},
	}
	return rsp
}

func TestLinkPostPostgres(t *testing.T) {
	SERVER_URL := os.Getenv("SERVER_URL")

	db := models.NewDatabase(utils.PostgresDB)
	if err := db.Connect(); err != nil {
		log.Fatalln(err)
	}
	err := db.Pst.QueryRow(context.Background(), "DELETE FROM link").Scan()
	if err != nil && err != pgx.ErrNoRows {
		t.Error(err)
	}
	svr := httptest.NewServer(controllers.NewLink(db, SERVER_URL))
	defer svr.Close()
	client := http.Client{}
	data := NewRspWithDataPost()
	for _, datum := range data {
		req, err := http.NewRequest("POST", svr.URL, bytes.NewBuffer(datum.data))
		if err != nil {
			t.Error(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, resp.StatusCode, datum.rspCode)
	}
}

func TestLinkGetPostgres(t *testing.T) {
	SERVER_URL := os.Getenv("SERVER_URL")

	db := models.NewDatabase(utils.PostgresDB)
	if err := db.Connect(); err != nil {
		log.Fatalln(err)
	}
	err := db.Pst.QueryRow(context.Background(), "DELETE FROM link").Scan()
	if err != nil && err != pgx.ErrNoRows {
		t.Error(err)
	}
	svr := httptest.NewServer(controllers.NewLink(db, SERVER_URL))
	defer svr.Close()
	client := http.Client{}
	data := NewRspWithDataPost()
	dataGet := make([]Rsp, 0, len(data))

	for _, datum := range data {
		req, err := http.NewRequest("POST", svr.URL, bytes.NewBuffer(datum.data))
		if err != nil {
			t.Error(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		url := services.URL{}
		err = json.NewDecoder(resp.Body).Decode(&url)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, resp.StatusCode, datum.rspCode)
		dataGet = append(dataGet, Rsp{rspCode: resp.StatusCode, url: datum.url, rspUrl: url.Url})

	}

	for _, rsp := range dataGet {
		if rsp.rspCode != http.StatusOK && rsp.rspCode != http.StatusCreated {
			continue
		}

		req, err := http.NewRequest("GET", svr.URL, bytes.NewBuffer([]byte(`{"url": `+`"`+rsp.rspUrl+`"}`)))
		if err != nil {
			t.Error(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		url := services.URL{}
		err = json.NewDecoder(resp.Body).Decode(&url)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, rsp.url, url.Url)
	}
}

var methods = []string{
	"PUT",
	"PATCH",
	"DELETE",
	"HEAD",
	"PATCH",
	"LINK",
	"UNLOCK",
	"COPY",
}

func TestLinkOtherMethods(t *testing.T) {
	SERVER_URL := os.Getenv("SERVER_URL")

	db := models.NewDatabase(utils.PostgresDB)
	if err := db.Connect(); err != nil {
		log.Fatalln(err)
	}
	err := db.Pst.QueryRow(context.Background(), "DELETE FROM link").Scan()
	if err != nil && err != pgx.ErrNoRows {
		t.Error(err)
	}
	svr := httptest.NewServer(controllers.NewLink(db, SERVER_URL))
	defer svr.Close()
	client := http.Client{}
	for _, m := range methods {
		req, err := http.NewRequest(m, svr.URL, nil)
		if err != nil {
			t.Error(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, resp.StatusCode, http.StatusMethodNotAllowed)
	}
}
