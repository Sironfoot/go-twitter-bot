package db

import (
	"database/sql"

	_ "github.com/lib/pq" // initialise postgresql DB provider
)

var db *sql.DB

// InitDB initialises the database
func InitDB(connectionString string) (err error) {
	db, err = sql.Open("postgres", connectionString)
	return
}

// Entity represents interface that all database mapped structs implement
type Entity interface {
	ID() string
	IsTransient() bool
	Save() error
	Delete() error
}
