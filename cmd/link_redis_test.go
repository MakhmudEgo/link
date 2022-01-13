package main

import (
	"bytes"
	"encoding/json"
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

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}

func TestHomeRedis(t *testing.T) {
	SERVER_URL := os.Getenv("SERVER_URL")

	db := models.NewDatabase(utils.RedisDB)
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

func TestLinkPostRedis(t *testing.T) {
	SERVER_URL := os.Getenv("SERVER_URL")

	db := models.NewDatabase(utils.RedisDB)
	if err := db.Connect(); err != nil {
		log.Fatalln(err)
	}
	db.Rdb.FlushAll(db.Ctx)

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

func TestLinkGetRedis(t *testing.T) {
	SERVER_URL := os.Getenv("SERVER_URL")

	db := models.NewDatabase(utils.RedisDB)
	if err := db.Connect(); err != nil {
		log.Fatalln(err)
	}
	db.Rdb.FlushAll(db.Ctx)

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
