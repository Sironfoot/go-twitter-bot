package db

import (
	"database/sql"
	"errors"
	"regexp"

	_ "github.com/lib/pq" // initialise postgresql DB provider
)

var db *sql.DB

// InitDB initialises the database
func InitDB(connectionString string) (err error) {
	db, err = sql.Open("postgres", connectionString)
	return
}

// CloseDB closes the database
func CloseDB() error {
	return db.Close()
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
	Limit    int
	OrderBy  string
	OrderAsc bool
	After    interface{}
}

var isUUID = regexp.MustCompile(`(?i)^[a-f0-9]{8}\-[a-f0-9]{4}\-[a-f0-9]{4}\-[a-f0-9]{4}\-[a-f0-9]{12}$`)
