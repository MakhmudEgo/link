package utils

import (
	"errors"
	"os"
)

const (
	PostgresDB = iota
	RedisDB
	BadDB
)

type IParse interface {
	Args() (int, error)
}

type Parse struct {
}

func (p Parse) Args() (int, error) {
	if len(os.Args) != 2 {
		return BadDB, errors.New("to many arguments")
	} else if os.Args[1] == "redis" {
		return RedisDB, nil
	} else if os.Args[1] == "postgres" {
		return PostgresDB, nil
	}
	return BadDB, errors.New("bad argument")
}
