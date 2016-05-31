package sqlboiler

import (
	"database/sql"
	"errors"
)

// Entity represents interface that all database mapped structs implement
type Entity interface {
	IsTransient() bool
	MetaData() EntityMetaData
}

// EntityMetaData provides mapping meta data about a database struct entity
type EntityMetaData struct {
	TableName      string
	PrimaryKeyName string
}

// DataAccessor defines an interface for various data access methods required by sqlboiler.
type DataAccessor interface {
	Prepare(sql string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// ErrEntityNotFound is returned when a database Entity is not found, returned
// from functions that return a single Entity (e.g. EntityFromID)
var ErrEntityNotFound = errors.New("sqlboiler: Entity not found")
