package sqlboiler

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// EntityGetByID returns an entity by its ID
func EntityGetByID(entity Entity, id interface{}, db *sql.DB) error {
	metaData := entity.MetaData()
	selectSQL := GenerateGetByIDStatement(entity)

	var fields []interface{}
	var pkFieldValue reflect.Value

	val := reflect.ValueOf(entity).Elem()
	entityType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldInfo := entityType.Field(i)
		tag := fieldInfo.Tag
		columnName := strings.TrimSpace(tag.Get("db"))

		if columnName != "" && columnName != metaData.PrimaryKeyName {
			field := val.Field(i)
			fields = append(fields, field.Addr().Interface())
		}

		if columnName == metaData.PrimaryKeyName {
			pkFieldValue = val.Field(i)
		}
	}

	if reflect.TypeOf(id) != pkFieldValue.Type() {
		return fmt.Errorf("id argument type (%s) doesn't match Entity ID type (%s)",
			reflect.TypeOf(id).Name(), pkFieldValue.Kind().String())
	}

	err := db.QueryRow(selectSQL, id).
		Scan(fields...)
	if err == sql.ErrNoRows {
		return ErrEntityNotFound
	} else if err != nil {
		return err
	}

	switch pkFieldValue.Kind() {
	case reflect.String:
		pkFieldValue.SetString(id.(string))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		pkFieldValue.SetInt(id.(int64))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		pkFieldValue.SetUint(id.(uint64))
	default:
		return fmt.Errorf("Entity (%s) primary key ID is a type (%s) that is not supported",
			entityType.Name(), pkFieldValue.Kind().String())
	}

	return nil
}

// EntitySave saves (either INSERTs or UPDATEs) an entity to the database
func EntitySave(entity Entity, db *sql.DB) error {
	metaData := entity.MetaData()

	var fields []interface{}
	var pkFieldValue reflect.Value

	val := reflect.ValueOf(entity).Elem()
	entityType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldInfo := entityType.Field(i)
		tag := fieldInfo.Tag
		columnName := strings.TrimSpace(tag.Get("db"))

		if columnName != "" && columnName != metaData.PrimaryKeyName {
			field := val.Field(i)
			fields = append(fields, field.Interface())
		}

		if columnName == metaData.PrimaryKeyName {
			pkFieldValue = val.Field(i)
		}
	}

	if entity.IsTransient() {
		insertSQL := GenerateInsertStatement(entity)

		statement, err := db.Prepare(insertSQL)
		if err != nil {
			return err
		}
		defer statement.Close()

		row := statement.QueryRow(fields...)

		switch pkFieldValue.Kind() {
		case reflect.String:
			var id string
			err = row.Scan(&id)
			if err == nil {
				pkFieldValue.SetString(id)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var id int64
			err = row.Scan(&id)
			if err == nil {
				pkFieldValue.SetInt(id)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var id uint64
			err = row.Scan(&id)
			if err == nil {
				pkFieldValue.SetUint(id)
			}
		default:
			err = fmt.Errorf("Entity (%s) primary key ID is a type (%s) that is not supported.",
				entityType.Name(), pkFieldValue.Kind().String())
		}

		if err != nil {
			return err
		}
	} else {
		updateSQL := GenerateUpdateStatement(entity)
		fields = append([]interface{}{pkFieldValue.Interface()}, fields...)

		_, err := db.Exec(updateSQL, fields...)
		if err != nil {
			return err
		}
	}

	return nil
}

// EntityDelete deletes an entity from the databse
func EntityDelete(entity Entity, db *sql.DB) error {
	metaData := entity.MetaData()

	var pkFieldValue reflect.Value

	val := reflect.ValueOf(entity).Elem()
	entityType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldInfo := entityType.Field(i)
		tag := fieldInfo.Tag
		columnName := strings.TrimSpace(tag.Get("db"))

		if columnName == metaData.PrimaryKeyName {
			pkFieldValue = val.Field(i)
		}
	}

	deleteSQL := GenerateDeleteByIDStatement(entity)
	_, err := db.Exec(deleteSQL, pkFieldValue.Interface())
	return err
}
