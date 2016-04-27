package sqlboiler

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// GetColumnList returns a slice of all the column names for a DB
// Entity struct (as specified by their tags), except the PK ID
func GetColumnList(entity Entity, alias string) []string {
	metaData := entity.MetaData()

	var columnList []string

	val := reflect.ValueOf(entity).Elem()
	entityType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldInfo := entityType.Field(i)
		tag := fieldInfo.Tag
		columnName := strings.TrimSpace(tag.Get("db"))

		if columnName != "" && columnName != metaData.PrimaryKeyName {
			if alias != "" {
				columnName = alias + "." + columnName
			}

			columnList = append(columnList, columnName)
		}
	}

	return columnList
}

// GetColumnListString returns a comma separated string of all the column
// names for a DB Entity struct (as specified by thier tags), except the PK ID
func GetColumnListString(entity Entity, alias string) string {
	return strings.Join(GetColumnList(entity, alias), ", ")
}

// GenerateInsertStatement generates an SQL command for inserting a record into the database
func GenerateInsertStatement(entity Entity) string {
	metaData := entity.MetaData()

	columns := GetColumnList(entity, "")

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

	columns := GetColumnList(entity, "")
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

	return "SELECT " + GetColumnListString(entity, "") + " " +
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
