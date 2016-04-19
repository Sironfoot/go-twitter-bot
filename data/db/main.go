package db

import (
	"database/sql"
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

// QueryAll is a general purpose query struct for returning Entities
type QueryAll struct {
	Limit    int
	OrderBy  string
	OrderAsc bool
	After    interface{}
}

// PagingInfo contains information about paging when calling queries that return multiple records
type PagingInfo struct {
	Limit   int
	Offset  int
	OrderBy string
	Asc     bool
}

var isUUID = regexp.MustCompile(`(?i)^[a-f0-9]{8}\-[a-f0-9]{4}\-[a-f0-9]{4}\-[a-f0-9]{4}\-[a-f0-9]{12}$`)
