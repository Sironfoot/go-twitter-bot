package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// EntityMetaData provides mapping meta data about a database struct entity
type EntityMetaData struct {
	TableName      string
	PrimaryKeyName string
}

// GetColumnList returns a slice of all the column names for a DB
// Entity struct (as specified by thier tags), except the PK ID
func GetColumnList(entity Entity) []string {
	metaData := entity.MetaData()

	var columnList []string

	val := reflect.ValueOf(entity).Elem()
	entityType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldInfo := entityType.Field(i)
		tag := fieldInfo.Tag
		columnName := strings.TrimSpace(tag.Get("db"))

		if columnName != "" && columnName != metaData.PrimaryKeyName {
			columnList = append(columnList, columnName)
		}
	}

	return columnList
}

// GetColumnListString returns a comma separated string of all the column
// names for a DB Entity struct (as specified by thier tags), except the PK ID
func GetColumnListString(entity Entity) string {
	return strings.Join(GetColumnList(entity), ", ")
}

// GenerateInsertStatement generates an SQL command for inserting a record into the database
func GenerateInsertStatement(entity Entity) string {
	metaData := entity.MetaData()

	columns := GetColumnList(entity)

	cmd := "INSERT INTO " + metaData.TableName + "(" + strings.Join(columns, ", ") + ") " +
		"VALUES("

	for i := range columns {
		cmd += fmt.Sprintf("$%d, ", i+1)
	}

	return cmd[:len(cmd)-2] + ") " +
		"RETURNING " + metaData.PrimaryKeyName
}

// GenerateUpdateStatement generates an SQL command for updating a record to the database
func GenerateUpdateStatement(entity Entity) string {
	metaData := entity.MetaData()

	cmd := "UPDATE " + metaData.TableName + " SET "

	columns := GetColumnList(entity)
	for i, columnName := range columns {
		cmd += fmt.Sprintf("%s = $%d, ", columnName, i+2)
	}

	cmd = cmd[:len(cmd)-2]
	return cmd + " WHERE " + metaData.PrimaryKeyName + " = $1"
}

// GenerateDeleteByIDStatement generates an SQL command to DELETE a record from the database by its ID
func GenerateDeleteByIDStatement(entity Entity) string {
	metaData := entity.MetaData()
	return "DELETE FROM " + metaData.TableName + " WHERE " + metaData.PrimaryKeyName + " = $1"
}

// GenerateGetByIDStatement generates an SQL command to SELECT the record from the database by its ID
func GenerateGetByIDStatement(entity Entity) string {
	metaData := entity.MetaData()

	return "SELECT " + GetColumnListString(entity) + " " +
		"FROM " + metaData.TableName + " " +
		"WHERE " + metaData.PrimaryKeyName + " = $1"
}

// GenerateGetAllStatement generates an SQL command to SELECT
// all records from the database complete with paging information
func GenerateGetAllStatement(entity Entity, where string) string {
	metaData := entity.MetaData()

	if strings.TrimSpace(where) != "" {
		// convert "email = $1" to "email = $4"
		where = regexp.MustCompile(`\$[0-9]+`).
			ReplaceAllStringFunc(where, func(placeholder string) string {
				placeholderNum, _ := strconv.Atoi(regexp.MustCompile(`[0-9]+`).FindString(placeholder))
				placeholderNum = placeholderNum + 3
				return "$" + strconv.Itoa(placeholderNum)
			})
	}

	var columnNames []string

	val := reflect.ValueOf(entity).Elem()
	entityType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldInfo := entityType.Field(i)
		tag := fieldInfo.Tag
		columnName := strings.TrimSpace(tag.Get("db"))

		if columnName != "" && columnName != metaData.PrimaryKeyName {
			columnNames = append(columnNames, columnName)
		}
	}

	cmd := "SELECT " + metaData.PrimaryKeyName + ", " + strings.Join(columnNames, ", ") + " " +
		"FROM " + metaData.TableName + " "

	if strings.TrimSpace(where) != "" {
		cmd += "WHERE " + where + " "
	}

	cmd += "ORDER BY $1 " +
		"LIMIT $2 OFFSET $3"

	return cmd
}

// EntityGetByID returns an entity by its ID
func EntityGetByID(entity Entity, id interface{}) error {
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
func EntitySave(entity Entity) error {
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
func EntityDelete(entity Entity) error {
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
