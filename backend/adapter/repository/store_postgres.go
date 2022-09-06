package repository

import (
	"database/sql"
	"github.com/jackc/pgx/v4"
)

type StorePostgreSQL struct {
	db *pgx.Conn
}

func NewStorePostgreSQL(db *pgx.Conn) *StorePostgreSQL {
	return &StorePostgreSQL{db: db}
}
