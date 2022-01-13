package main

import (
	"github.com/joho/godotenv"
	"link/internal/pkg/utils"
	"link/pkg/controllers"
	"link/pkg/models"
	"log"
	"net/http"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}

func main() {
	SERVER_URL := os.Getenv("SERVER_URL")
	parse := utils.Parse{}
	typeDB, err := parse.Args()
	if err != nil {
		log.Fatalln(err)
	}
	db := models.NewDatabase(typeDB)
	if err := db.Connect(); err != nil {
		log.Fatalln(err)
	}
	if typeDB == utils.PostgresDB {
		defer db.Close()
	}

	// endpoints
	http.Handle("/", controllers.NewHome(db, SERVER_URL))
	http.Handle("/link", controllers.NewLink(db, SERVER_URL))

	log.Fatalln(http.ListenAndServe(os.Getenv("SERVER_PORT"), nil))
}
