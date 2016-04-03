package models

import (
	"database/sql"
	"io/ioutil"
	"testing"

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

	createTables := "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"

	file, err := ioutil.ReadFile("../sql/create-tables.sql")
	if err != nil {
		return err
	}
	createTables += string(file)

	_, err = testDB.Exec(createTables)
	if err != nil {
		return err
	}

	return nil
}

func tearDown() error {
	tempDB, err := sql.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		return err
	}

	_, err = tempDB.Exec("DROP DATABASE IF EXISTS go_twitter_bot_test;")
	if err != nil {
		return err
	}

	tempDB.Close()
	testDB.Close()

	return nil
}

func TestUserValidateCreate(t *testing.T) {
	err := setUp()
	if err != nil {
		t.Fatal(err)
	}

	err = db.InitDB("user=postgres dbname=go_twitter_bot_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	user := User{
		Email:    "someone@example.com",
		Password: "Password1",
		IsAdmin:  true,
	}

	validationErrors, err := user.ValidateCreate()
	if err != nil {
		t.Fatal(err)
	}

	if len(validationErrors) > 0 {
		t.Errorf("should have passed validation, but got %d errors.", len(validationErrors))
	}

	err = tearDown()
	if err != nil {
		t.Fatal(err)
	}
}
