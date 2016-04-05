package db_test

import (
	"database/sql"
	"io/ioutil"
	"testing"
	"time"

	"github.com/sironfoot/go-twitter-bot/data/db"
)

var testDB *sql.DB

func setUp() (err error) {
	tempDB, err := sql.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		return err
	}

	_, err = tempDB.Exec("DROP DATABASE IF EXISTS go_twitter_bot_test;")
	if err != nil {
		return err
	}

	_, err = tempDB.Exec("CREATE DATABASE go_twitter_bot_test;")
	if err != nil {
		return err
	}

	tempDB.Close()

	testDB, err = sql.Open("postgres", "user=postgres dbname=go_twitter_bot_test sslmode=disable")
	if err != nil {
		return err
	}

	createTables := "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";\n"

	file, err := ioutil.ReadFile("../sql/create-tables.sql")
	if err != nil {
		return err
	}
	createTables += string(file)

	_, err = testDB.Exec(createTables)
	if err != nil {
		return err
	}

	return db.InitDB("user=postgres dbname=go_twitter_bot_test sslmode=disable")
}

func tearDown() error {
	if err := db.CloseDB(); err != nil {
		return err
	}

	if err := testDB.Close(); err != nil {
		return err
	}

	tempDB, err := sql.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		return err
	}

	_, err = tempDB.Exec("DROP DATABASE IF EXISTS go_twitter_bot_test;")
	if err != nil {
		return err
	}

	if err := tempDB.Close(); err != nil {
		return err
	}

	return nil
}

func createTestUser() (user db.User, id string, err error) {
	user = db.User{
		Email:          "test@example.com",
		HashedPassword: "Password1",
		IsAdmin:        true,
		DateCreated:    time.Now(),
	}

	createSQL := `INSERT INTO users(email, hashed_password, is_admin, date_created)
                  VALUES($1, $2, $3, $4)
                  RETURNING id`

	statement, err := testDB.Prepare(createSQL)
	if err != nil {
		return
	}
	defer statement.Close()

	err = statement.
		QueryRow(user.Email, user.HashedPassword, user.IsAdmin, user.DateCreated).
		Scan(&id)
	if err != nil {
		return
	}

	return
}

func TestUserFromID(t *testing.T) {
	if err := setUp(); err != nil {
		t.Fatal(err)
		return
	}

	defer func() {
		if err := tearDown(); err != nil {
			t.Fatal(err)
		}
	}()

	// arrange
	testUser, id, err := createTestUser()

	// Non-existent record - Valid UUID
	// act
	user, err := db.UserFromID("ec24d9b2-fb39-11e5-9dcc-df8c5db12101")
	if err != nil && err != db.ErrEntityNotFound {
		t.Fatal(err)
	}

	// assert
	if err != db.ErrEntityNotFound {
		t.Errorf("user entity was returned from non-existent ID, userID: %s", user.ID())
	}

	// Non-existent record - Invalid UUID
	user, err = db.UserFromID("Nonsense")
	if err != nil && err != db.ErrEntityNotFound {
		t.Fatal(err)
	}

	// assert
	if err != db.ErrEntityNotFound {
		t.Errorf("user entity was returned from non-existent invalid ID, userID: %s", user.ID())
	}

	// Existing record
	// act
	user, err = db.UserFromID(id)
	if err != nil && err != db.ErrEntityNotFound {
		t.Fatal(err)
	}

	// assert
	if err == db.ErrEntityNotFound {
		t.Errorf("Expected user record, but got ErrEntityNotFound")
	}

	if user.ID() != id ||
		user.Email != testUser.Email ||
		user.IsAdmin != testUser.IsAdmin ||
		user.DateCreated.Equal(testUser.DateCreated) {

		t.Errorf("Expected user and actual user don't match")
	}
}

func TestUserFromEmail(t *testing.T) {
	if err := setUp(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := tearDown(); err != nil {
			t.Fatal(err)
		}
	}()

	// arrange
	testUser, id, err := createTestUser()

	// Non-existent record
	// act
	user, err := db.UserFromEmail("Nonsense@example.com")
	if err != nil && err != db.ErrEntityNotFound {
		t.Fatal(err)
	}

	// assert
	if err != db.ErrEntityNotFound {
		t.Errorf("user entity was returned from non-existent email address, userEmail: %s", user.Email)
	}

	// Existing record
	// act
	user, err = db.UserFromEmail(testUser.Email)
	if err != nil && err != db.ErrEntityNotFound {
		t.Fatal(err)
	}

	// assert
	if err == db.ErrEntityNotFound {
		t.Errorf("Expected user record, but got ErrEntityNotFound")
	}

	if user.ID() != id ||
		user.Email != testUser.Email ||
		user.IsAdmin != testUser.IsAdmin ||
		user.DateCreated.Equal(testUser.DateCreated) {

		t.Errorf("Expected user and actual user don't match")
	}
}
