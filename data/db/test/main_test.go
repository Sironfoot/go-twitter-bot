package db_test

import (
	"database/sql"
	"io/ioutil"
	"time"

	"github.com/sironfoot/go-twitter-bot/data/db"
)

var testDB *sql.DB

func mustSetUp() {
	tempDB, err := sql.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	_, err = tempDB.Exec("DROP DATABASE IF EXISTS go_twitter_bot_test;")
	if err != nil {
		panic(err)
	}

	_, err = tempDB.Exec("CREATE DATABASE go_twitter_bot_test;")
	if err != nil {
		panic(err)
	}

	tempDB.Close()

	testDB, err = sql.Open("postgres", "user=postgres dbname=go_twitter_bot_test sslmode=disable")
	if err != nil {
		panic(err)
	}

	createTables := "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";\n"

	file, err := ioutil.ReadFile("../../sql/create-tables.sql")
	if err != nil {
		panic(err)
	}
	createTables += string(file)

	_, err = testDB.Exec(createTables)
	if err != nil {
		panic(err)
	}

	err = db.InitDB("user=postgres dbname=go_twitter_bot_test sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func mustTearDown() {
	if err := db.CloseDB(); err != nil {
		panic(err)
	}

	if err := testDB.Close(); err != nil {
		panic(err)
	}

	tempDB, err := sql.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	_, err = tempDB.Exec("DROP DATABASE IF EXISTS go_twitter_bot_test;")
	if err != nil {
		panic(err)
	}

	if err := tempDB.Close(); err != nil {
		panic(err)
	}

	return
}

func createTestUser() (user db.User, id string, err error) {
	return createTestUserWithEmail("test@example.com")
}

func createTestUserWithEmail(email string) (user db.User, id string, err error) {
	user = db.User{
		Email:          email,
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
