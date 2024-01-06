package gift

import (
	"context"
	"database/sql"
	"discount/db"
	"errors"
	"github.com/redis/go-redis/v9"
)

var (
	ErrNoRowToUpdate = errors.New("no row to update")
)

type Storage struct {
	db    db.SQLExt
	redis *redis.Client
}

func New(db *sql.DB, redis *redis.Client) Storage {
	println("\033[31m" + "gift storage redis" + redis.Ping(context.Background()).String() + "\033[0m")
	return Storage{db: db, redis: redis}
}

func (s Storage) WithTX(tx *sql.Tx) (Storage, error) {
	if tx == nil {
		return Storage{}, db.ErrNoTXProvided
	}
	switch s.db.(type) {
	case *sql.Tx:
		return Storage{}, db.ErrAlreadyInTX
	case *sql.DB:
		return Storage{db: tx}, nil
	}
	return s, nil
}
