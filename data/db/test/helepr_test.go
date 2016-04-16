package db_test

import (
	"testing"

	"github.com/sironfoot/go-twitter-bot/data/db"
)

func TestGenerateInsertStatement(t *testing.T) {
	user := db.User{}

	expected := "INSERT INTO users(name, email, hashed_password, auth_token, is_admin, is_service, date_created) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7) " +
		"RETURNING id"
	actual := db.GenerateInsertStatement(&user)

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}
}

func TestGenerateUpdateStatement(t *testing.T) {
	user := db.User{}

	expected := "UPDATE users " +
		"SET name = $2, " +
		"email = $3, " +
		"hashed_password = $4, " +
		"auth_token = $5, " +
		"is_admin = $6, " +
		"is_service = $7, " +
		"date_created = $8 " +
		"WHERE id = $1"

	actual := db.GenerateUpdateStatement(&user)

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}
}

func TestGenerateDeleteByIDStatement(t *testing.T) {
	user := db.User{}

	expected := "DELETE FROM users WHERE id = $1"
	actual := db.GenerateDeleteByIDStatement(&user)

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}
}

func TestGenerateGetByIdStatement(t *testing.T) {
	user := db.User{}

	expected := "SELECT * FROM users WHERE id = $1"
	actual := db.GenerateGetByIDStatement(&user)

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}
}

func TestGenerateGetAllStatement(t *testing.T) {
	user := db.User{}

	// no WHERE clause
	expected := "SELECT * " +
		"FROM users " +
		"ORDER BY $1 " +
		"LIMIT $2 OFFSET $3"

	actual := db.GenerateGetAllStatement(&user, "")

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}

	// with WHERE clause
	expected = "SELECT * " +
		"FROM users " +
		"WHERE email = $4 AND name LIKE '$5' " +
		"ORDER BY $1 " +
		"LIMIT $2 OFFSET $3"
	actual = db.GenerateGetAllStatement(&user, "email = $1 AND name LIKE '$2'")

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}
}
