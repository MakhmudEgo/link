package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/cmd/controllers"
	"main/cmd/models"
	"main/cmd/utils"
	"net/http"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}

func main() {
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
	http.Handle("/", controllers.NewHome(db))
	http.Handle("/link", controllers.NewLink(db))

	log.Fatalln(http.ListenAndServe(os.Getenv("SERVER_PORT"), nil))
}
