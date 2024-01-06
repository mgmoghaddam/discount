package discount

import (
	"database/sql"
	"discount/db"
	"errors"
)

var (
	ErrNoRowToUpdate = errors.New("no row to update")
)

type Storage struct {
	db db.SQLExt
}

func New(db *sql.DB) Storage {
	return Storage{db: db}
}

// WithTX returns a new storage with the given transaction replacing the db.
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
