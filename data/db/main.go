package db

import (
	"database/sql"
	"errors"
	"regexp"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // initialise postgresql DB provider
)

var db *sql.DB
var dbx *sqlx.DB

// InitDB initialises the database
func InitDB(connectionString string) error {
	var err error

	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	dbx = sqlx.NewDb(db, "postgres")
	return nil
}

// CloseDB closes the database
func CloseDB() error {
	err := db.Close()
	if err != nil {
		return err
	}

	return dbx.Close()
}

// PagingInfo contains information about paging when calling queries that return multiple records
type PagingInfo struct {
	Limit   int
	Offset  int
	OrderBy string
	Asc     bool
}

// ErrEntityNotFound is returned when a database Entity is not found, returned
// from functions that return a single Entity (e.g. EntityFromID)
var ErrEntityNotFound = errors.New("db: Entity not found")

var isUUID = regexp.MustCompile(`(?i)^[a-f0-9]{8}\-[a-f0-9]{4}\-[a-f0-9]{4}\-[a-f0-9]{4}\-[a-f0-9]{12}$`)
