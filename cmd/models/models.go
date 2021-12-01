package models

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"log"
	"main/cmd/utils"
	"os"
)

type IModels interface {
	Select() error
	Connect() error
	Create()
	Check()
}

type Link struct {
	Base  string `json:"base"`
	Short string `json:"short"`
}

type Database struct {
	typeBD int
	rdb    *redis.Client
	pst    *pgx.Conn
	ctx    context.Context
}

func NewDatabase(typeBD int) *Database {
	return &Database{typeBD: typeBD, ctx: context.Background()}
}

func (db *Database) Connect() error {
	if err := db.Select(); err != nil {
		return err
	}
	return nil
}

func (db *Database) Select() error {
	var err error
	switch db.typeBD {
	case utils.PostgresDB:
		db.pst, err = pgx.Connect(db.ctx, os.Getenv("DATABASE_URL"))
		if err != nil {
			return err
		} else {
			log.Println("postgres success")
		}
	case utils.RedisDB:
		db.rdb = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
		})
		err = db.rdb.Set(db.ctx, "key", "value", 0).Err()
		if err != nil {
			return err
		} else {
			log.Println("redis success")
		}

	}
	return nil
}
