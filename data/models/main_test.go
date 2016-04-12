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

type testCase struct {
	model          Model
	id             string
	expectedErrors map[string]string
}

func runValidationTest(t *testing.T, testCases []testCase, runValidation func(model Model, id string) ([]ValidationError, error)) {
	for i, testCase := range testCases {
		validationErrors, err := runValidation(testCase.model, testCase.id)
		if err != nil {
			t.Fatal(err)
		}

		if len(validationErrors) != len(testCase.expectedErrors) {
			t.Errorf("test case %d: expected %d validation error(s) but got %d",
				i+1, len(testCase.expectedErrors), len(validationErrors))
		} else {
			for _, validationError := range validationErrors {
				expectedType, ok := testCase.expectedErrors[validationError.FieldName]

				if !ok {
					t.Errorf("test case %d: validationError field '%s' wasn't expected",
						i+1, validationError.FieldName)
				}

				if ok && validationError.Type != expectedType {
					t.Errorf("test case %d: validationError field '%s' type (%s) doesn't match expected type (%s)",
						i+1, validationError.FieldName, validationError.Type, expectedType)
				}
			}
		}
	}
}
