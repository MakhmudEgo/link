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
	db         *models.Database
	SERVER_URL string
}

type Link struct {
	db         *models.Database
	SERVER_URL string
}

func NewHome(db *models.Database, SERVER_URL string) *Home {
	return &Home{db: db, SERVER_URL: SERVER_URL}
}

func NewLink(db *models.Database, SERVER_URL string) *Link {
	return &Link{db: db, SERVER_URL: SERVER_URL}
}

func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte(`<h1 style="
			display: flex;
			align-items: center; 
			justify-content: center; 
			height: 100%; 
			font-size: 150px;
		">☕️</h1>`))
	//http.Redirect(w, r, "https://yandex.ru", http.StatusTemporaryRedirect)

}

func (l *Link) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	service := services.NewLink(l.db, r, l.SERVER_URL)
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
