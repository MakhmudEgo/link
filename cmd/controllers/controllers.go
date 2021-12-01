package controllers

import (
	"encoding/json"
	"log"
	"main/cmd/models"
	"main/cmd/services"
	"net/http"
)

/*type IControllers interface {
	Home()
}*/

type Home struct {
	db *models.Database
}

func (h *Home) name() {

}

func NewHome(db *models.Database) *Home {
	return &Home{db: db}
}

type Link struct {
	db *models.Database
}

func NewLink(db *models.Database) *Link {
	return &Link{db: db}
}

func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//http.Redirect(w, r, "https://yandex.ru", http.StatusTemporaryRedirect)

}

func (l *Link) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	service := services.NewLink(l.db, r)
	resp, err := service.LinkExecute()
	js, errJson := json.Marshal(resp)
	if errJson != nil {
		log.Fatalln(errJson)
	}

	w.WriteHeader(resp.Code)

	if _, errWrite := w.Write(js); err != nil {
		log.Fatalln(errWrite)
	}

	if err != nil {
		log.Fatalln(err)
	}

}
