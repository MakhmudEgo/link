package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/cmd/controllers"
	"main/cmd/models"
	"main/cmd/utils"
	"net/http"
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

	http.Handle("/", &controllers.Home{})
	http.Handle("/link", &controllers.Link{})

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
