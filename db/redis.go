package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(host, password, port string, db int, timeout time.Duration) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%s", host, port),
		Password:    password,
		DB:          db,
		ReadTimeout: timeout,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return rdb, nil
}
