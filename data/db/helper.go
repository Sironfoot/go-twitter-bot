package db

import (
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

// GenerateInsertStatement generates an SQL command for inserting a record into the database
func GenerateInsertStatement(entity Entity) string {
	metaData := entity.MetaData()

	insertLine := "INSERT INTO " + metaData.TableName + "("
	valuesLine := "VALUES("

	val := reflect.ValueOf(entity).Elem()
	entityType := val.Type()

	placeholder := 1

	for i := 0; i < val.NumField(); i++ {
		fieldInfo := entityType.Field(i)
		tag := fieldInfo.Tag
		columnName := strings.TrimSpace(tag.Get("db"))

		if columnName != "" && columnName != metaData.PrimaryKeyName {
			insertLine += columnName + ", "
			valuesLine += fmt.Sprintf("$%d, ", placeholder)
			placeholder++
		}
	}

	insertLine = insertLine[:len(insertLine)-2] + ") "
	valuesLine = valuesLine[:len(valuesLine)-2] + ") "

	return insertLine + valuesLine + "RETURNING " + metaData.PrimaryKeyName
}

// GenerateUpdateStatement generates an SQL command for updating a record to the database
func GenerateUpdateStatement(entity Entity) string {
	metaData := entity.MetaData()

	cmd := "UPDATE " + metaData.TableName + " SET "

	val := reflect.ValueOf(entity).Elem()
	entityType := val.Type()

	placeholder := 2

	for i := 0; i < val.NumField(); i++ {
		fieldInfo := entityType.Field(i)
		tag := fieldInfo.Tag
		columnName := strings.TrimSpace(tag.Get("db"))

		if columnName != "" && columnName != metaData.PrimaryKeyName {
			cmd += fmt.Sprintf("%s = $%d, ", columnName, placeholder)
			placeholder++
		}
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
	return "SELECT * FROM " + metaData.TableName + " WHERE " + metaData.PrimaryKeyName + " = $1"
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

	cmd := "SELECT * " +
		"FROM " + metaData.TableName + " "

	if strings.TrimSpace(where) != "" {
		cmd += "WHERE " + where + " "
	}

	cmd += "ORDER BY $1 " +
		"LIMIT $2 OFFSET $3"

	return cmd
}
