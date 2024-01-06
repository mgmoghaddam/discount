package main

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
)

func healthFunc(db *sql.DB, rdb *redis.Client) func() error {
	return func() error {
		if err := db.Ping(); err != nil {
			return err
		}
		if err := rdb.Ping(context.Background()).Err(); err != nil {
			return err
		}
		return nil
	}
}
