package models

import (
	"testing"
	"time"

	"github.com/sironfoot/go-twitter-bot/data/db"
)

var userFromEmail func(email string) (db.User, error)

func setUp() {
	userFromEmail = db.UserFromEmail
	db.UserFromEmail = func(email string) (db.User, error) {
		if email == "existing@example.com" {
			return db.User{
				Email:          email,
				HashedPassword: "Password1",
				IsAdmin:        true,
				DateCreated:    time.Now(),
			}, nil
		}

		return db.User{}, db.ErrEntityNotFound
	}
}

func tearDown() {
	db.UserFromEmail = userFromEmail
}

func TestUserValidateCreate(t *testing.T) {
	setUp()
	defer tearDown()

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
		t.Errorf("should have passed validation, but got %d error(s): %s", len(validationErrors), validationErrors)
	}
}
