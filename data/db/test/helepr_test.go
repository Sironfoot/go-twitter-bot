package db_test

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/sironfoot/go-twitter-bot/data/db"
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

func TestGenerateInsertStatement(t *testing.T) {
	user := db.User{}

	expected := "INSERT INTO users(" + strings.Join(columnNames, ", ") + ") " +
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

	expected := "SELECT " + strings.Join(columnNames, ", ") + " FROM users WHERE id = $1"
	actual := db.GenerateGetByIDStatement(&user)

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}
}

func TestGenerateGetAllStatement(t *testing.T) {
	user := db.User{}

	// no WHERE clause
	expected := "SELECT id, " + strings.Join(columnNames, ", ") + " " +
		"FROM users " +
		"ORDER BY $1 " +
		"LIMIT $2 OFFSET $3"

	actual := db.GenerateGetAllStatement(&user, "")

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}

	// with WHERE clause
	expected = "SELECT id, " + strings.Join(columnNames, ", ") + " " +
		"FROM users " +
		"WHERE email = $4 AND name LIKE '$5' " +
		"ORDER BY $1 " +
		"LIMIT $2 OFFSET $3"
	actual = db.GenerateGetAllStatement(&user, "email = $1 AND name LIKE '$2'")

	if expected != actual {
		t.Errorf("actual SQL statement was:\n\n%s\n\nbut expected:\n\n%s", actual, expected)
	}
}
func TestEntitySave(t *testing.T) {
	mustSetUp()
	defer mustTearDown()

	user := db.User{
		Name:           "Test User",
		Email:          "test@example.com",
		HashedPassword: "Password1_",
		AuthToken:      sql.NullString{String: "my_auth_token", Valid: true},
		IsAdmin:        true,
		IsService:      false,
		DateCreated:    time.Now().UTC(),
	}

	err := db.EntitySave(&user)
	if err != nil {
		t.Fatal(err)
	}

	if user.IsTransient() {
		t.Errorf("user is still in a transient state (non-saved), user ID: %s", user.ID)
	}

	user.Name = "Updated User"
	user.Email = "updated@example.com"
	user.HashedPassword = "UpdatedPassword1"
	user.AuthToken = sql.NullString{String: "updated_auth_token", Valid: true}
	user.IsAdmin = false
	user.IsService = true
	user.DateCreated = time.Now().UTC()

	previousID := user.ID

	err = db.EntitySave(&user)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != previousID {
		t.Errorf("user ID was updated from %s to %s", previousID, user.ID)
	}
}
