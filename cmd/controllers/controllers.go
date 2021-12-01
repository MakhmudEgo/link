package controllers

import (
	"context"
	"log"
	"net/http"
)

/*type IControllers interface {
	Home()
}*/

type Home struct {
	ctx context.Context
}

type Link struct {
}

func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ctx = context.Background()
}

func (l *Link) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"kek": "LOL"}`))
	if err != nil {
		log.Fatalln(err)
	}
	//http.Redirect(w, r, "https://yandex.ru", http.StatusTemporaryRedirect)
}
