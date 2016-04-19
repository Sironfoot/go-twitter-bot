package sqlboiler_test

import (
	"strings"
	"testing"

	"github.com/sironfoot/go-twitter-bot/lib/sqlboiler"
)

var columnNames = []string{
	"name",
	"email",
	"hashed_password",
	"auth_token",
	"is_admin",
	"is_service",
	"date_created",
}

func TestGenerateColumnList(t *testing.T) {
	user := userEntity{}

	expected := strings.Join(columnNames, ", ")
	actual := sqlboiler.GetColumnListString(&user)

	if expected != actual {
		t.Errorf("actual column list was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}
}

func TestGenerateInsertStatement(t *testing.T) {
	user := userEntity{}

	expected := "INSERT INTO users(" + strings.Join(columnNames, ", ") + ") " +
		"VALUES($1, $2, $3, $4, $5, $6, $7) " +
		"RETURNING id"
	actual := sqlboiler.GenerateInsertStatement(&user)

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}
}

func TestGenerateUpdateStatement(t *testing.T) {
	user := userEntity{}

	expected := "UPDATE users " +
		"SET name = $2, " +
		"email = $3, " +
		"hashed_password = $4, " +
		"auth_token = $5, " +
		"is_admin = $6, " +
		"is_service = $7, " +
		"date_created = $8 " +
		"WHERE id = $1"

	actual := sqlboiler.GenerateUpdateStatement(&user)

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}
}

func TestGenerateDeleteByIDStatement(t *testing.T) {
	user := userEntity{}

	expected := "DELETE FROM users WHERE id = $1"
	actual := sqlboiler.GenerateDeleteByIDStatement(&user)

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}
}

func TestGenerateGetByIdStatement(t *testing.T) {
	user := userEntity{}

	expected := "SELECT " + strings.Join(columnNames, ", ") + " FROM users WHERE id = $1"
	actual := sqlboiler.GenerateGetByIDStatement(&user)

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}
}

func TestGenerateGetAllStatement(t *testing.T) {
	user := userEntity{}

	// no WHERE clause
	expected := "SELECT id, " + strings.Join(columnNames, ", ") + " " +
		"FROM users " +
		"ORDER BY $1 " +
		"LIMIT $2 OFFSET $3"

	actual := sqlboiler.GenerateGetAllStatement(&user, "")

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}

	// with WHERE clause
	expected = "SELECT id, " + strings.Join(columnNames, ", ") + " " +
		"FROM users " +
		"WHERE email = $4 AND name LIKE '$5' " +
		"ORDER BY $1 " +
		"LIMIT $2 OFFSET $3"
	actual = sqlboiler.GenerateGetAllStatement(&user, "email = $1 AND name LIKE '$2'")

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}
}
