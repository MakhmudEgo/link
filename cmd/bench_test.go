package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/pkg/controllers"
	"main/pkg/models"
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

func BenchmarkHomePostgresGet(b *testing.B) {
	SERVER_URL := os.Getenv("SERVER_URL")

	db := models.NewDatabase(0)
	if err := db.Connect(); err != nil {
		log.Fatalln(err)
	}
	svr := httptest.NewServer(controllers.NewHome(db, SERVER_URL))
	defer svr.Close()
	req, err := http.NewRequest("GET", svr.URL, nil)
	if err != nil {
		b.Error(err)
	}
	client := http.Client{}
	for i := 0; i < b.N; i++ {
		_, err := client.Do(req)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkHomeRedisGet(b *testing.B) {
	SERVER_URL := os.Getenv("SERVER_URL")

	db := models.NewDatabase(1)
	if err := db.Connect(); err != nil {
		log.Fatalln(err)
	}
	svr := httptest.NewServer(controllers.NewHome(db, SERVER_URL))
	defer svr.Close()
	req, err := http.NewRequest("GET", svr.URL, nil)
	if err != nil {
		b.Error(err)
	}
	client := http.Client{}
	for i := 0; i < b.N; i++ {
		_, err := client.Do(req)
		if err != nil {
			b.Error(err)
		}
	}
}
