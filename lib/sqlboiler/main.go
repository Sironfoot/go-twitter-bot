package sqlboiler

import "errors"

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

// ErrEntityNotFound is returned when a database Entity is not found, returned
// from functions that return a single Entity (e.g. EntityFromID)
var ErrEntityNotFound = errors.New("sqlboiler: Entity not found")
