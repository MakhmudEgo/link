package services

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"log"
	"main/cmd/models"
	"main/cmd/utils"
	"net/http"
)

func NewHome(db *models.Database, r *http.Request, SERVER_URL string) *Home {
	return &Home{db: db, r: r, SERVER_URL: SERVER_URL}
}

func (h *Home) HomeExecute() (string, int, error) {
	res, err := h.Select(h.SERVER_URL + h.r.RequestURI[1:])
	if err != pgx.ErrNoRows && err != redis.Nil && err != nil {
		log.Println(err)
		return "", http.StatusInternalServerError, errors.New("Internal Server Error")
	} else if err == pgx.ErrNoRows || err == redis.Nil {
		log.Println(err)
		return "", http.StatusNotFound, errors.New("Not Found")
	}
	return res.Base, -1, nil
}

func (h *Home) Select(url string) (*models.Link, error) {
	log.Println("home.go, request url:", url)
	link := &models.Link{}
	var err error
	switch h.db.TypeDB {
	case utils.PostgresDB:
		err = h.db.Pst.QueryRow(
			context.Background(),
			"select base from link where short=$1",
			url).Scan(&link.Base)
	case utils.RedisDB:
		link.Base, err = h.db.Rdb.Get(h.db.Ctx, url).Result()
	}
	return link, err
}
