package main

import (
	"log"
	"main/cmd/controllers"
	"main/cmd/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHome(t *testing.T) {
	SERVER_URL := os.Getenv("SERVER_URL")

	db := models.NewDatabase(0)
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
	if resp.StatusCode != http.StatusNotFound {
		t.Error("StatusNotFound")
	}
	println(svr.URL)
	println(resp.StatusCode)

}

func TestLink(t *testing.T) {
	SERVER_URL := os.Getenv("SERVER_URL")

	db := models.NewDatabase(0)
	if err := db.Connect(); err != nil {
		log.Fatalln(err)
	}
	svr := httptest.NewServer(controllers.NewLink(db, SERVER_URL))
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
	if resp.StatusCode != http.StatusBadRequest {
		t.Error("StatusBadRequest")
	}
	println(svr.URL)
	println(resp.StatusCode)

}
