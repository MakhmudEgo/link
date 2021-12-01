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
	Close()
}

type Link struct {
	Base  string `json:"base"`
	Short string `json:"short"`
}

type Database struct {
	TypeDB int
	Rdb    *redis.Client
	Pst    *pgx.Conn
	Ctx    context.Context
}

func NewDatabase(typeBD int) *Database {
	return &Database{TypeDB: typeBD, Ctx: context.Background()}
}

func (db *Database) Connect() error {
	if err := db.Select(); err != nil {
		return err
	}
	return nil
}

func (db *Database) Select() error {
	var err error
	switch db.TypeDB {
	case utils.PostgresDB:
		db.Pst, err = pgx.Connect(db.Ctx, os.Getenv("DATABASE_URL"))
		if err != nil {
			return err
		} else {
			log.Println("postgres success")
		}
	case utils.RedisDB:
		db.Rdb = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
		})
		err = db.Rdb.Set(db.Ctx, "key", "value", 0).Err()
		if err != nil {
			return err
		} else {
			log.Println("redis success")
		}

	}
	return nil
}

func (db *Database) Close() {
	if db.TypeDB != utils.PostgresDB {
		return
	}
	if err := db.Pst.Close(db.Ctx); err != nil {
		log.Fatalln(err)
	}
}
