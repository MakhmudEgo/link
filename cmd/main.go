package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("to many args")
	} else if os.Args[1] == "redis" {

	} else if os.Args[1] == "postgres" {

	} else {
		log.Fatalln("bad arg")
	}
	http.HandleFunc("/link", func(w http.ResponseWriter, r *http.Request) {
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
	})

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
