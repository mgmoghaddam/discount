package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrAlreadyInTX    = errors.New("storage already running in a tx")
	ErrNoTXProvided   = errors.New("no tx provided")
	ErrDBNoTInitiated = errors.New("db not initiated")
)

var globalDB *sql.DB

func sqlDB() (*sql.DB, error) {
	if globalDB == nil {
		return nil, ErrDBNoTInitiated
	}
	return globalDB, nil
}

func NewPostgres(
	dbName, username, password, host, port string, maxOpenConnections, maxIdleConnections int,
) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, host, port, dbName,
	))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetMaxOpenConns(maxOpenConnections)
	if globalDB == nil {
		globalDB = db
	}
	return db, nil
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

type SQLExt interface {
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

func Transaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	db, err := sqlDB()
	if err != nil {
		return err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
