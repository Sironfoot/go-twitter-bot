package db

import (
	"database/sql"
	"errors"

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

// ErrEntityNotFound is returned when a database Entity is not found, returned
// from functions that return a single Entity (e.g. EntityFromID)
var ErrEntityNotFound = errors.New("db: Entity not found")

// QueryAll is a general purpose query struct for returning Entities
type QueryAll struct {
	startRecord int
	endRecord   int
	orderBy     string
	orderAsc    bool
}
