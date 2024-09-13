package model

import "database/sql"

type Model struct {
	conn *sql.DB
}

func New(db *sql.DB) *Model {
	return &Model{conn: db}
}
