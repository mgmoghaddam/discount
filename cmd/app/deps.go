package main

import (
	"database/sql"
	"discount/db"
	"discount/internal/config"
	"discount/server"
	"github.com/redis/go-redis/v9"
	"log"
)

func postgresDB() *sql.DB {
	psql, err := db.NewPostgres(
		config.DBName(), config.DBUser(), config.DBPassword(), config.DBHost(), config.DBPort(),
		config.DBMaxOpenConn(), config.DBMaxIdleConn(),
	)
	if err != nil {
		log.Fatalf("failed to initalize db: %v", err)
	}
	return psql
}

func redisDB() *redis.Client {
	rdb, err := db.NewRedis(config.RDBHost(), config.RDBPassword(), config.RDBPort(), config.RDB(), config.RDBTimeOut())
	if err != nil {
		log.Fatalf("failed to initalize redis: %v", err)
	}
	return rdb
}

func setupServer(s *server.Server, psql *sql.DB, rdb *redis.Client) {
	s.SetHealthFunc(healthFunc(psql, rdb)).
		SetupRoutes()
}
