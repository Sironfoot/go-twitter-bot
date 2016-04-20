package models_test

import (
	"testing"
	"time"

	"github.com/sironfoot/go-twitter-bot/data/db"
	"github.com/sironfoot/go-twitter-bot/data/models"
)

const (
	existingUserID    = "Test_ID"
	existingUserEmail = "existing@example.com"
)

var userFromEmail func(email string) (db.User, error)

func userSetUp() {
	userFromEmail = db.UserFromEmail
	db.UserFromEmail = func(email string) (db.User, error) {
		if email == existingUserEmail {
			return db.User{
				ID:             existingUserID,
				Name:           "Test Iser",
				Email:          existingUserEmail,
				HashedPassword: "Password1",
				IsAdmin:        true,
				IsService:      false,
				DateCreated:    time.Now(),
			}, nil
		}

		return db.User{}, db.ErrEntityNotFound
	}
}

func userTearDown() {
	db.UserFromEmail = userFromEmail
}

var commonTestCases = []testCase{
	{
		description: "no errors",
		model: &models.User{
			Name:     "Test User",
			Email:    "someone@example.com",
			Password: "Password1",
		},
		expectedErrors: []expectedError{},
	},
	{
		description: "email address invalid format",
		model: &models.User{
			Name:     "Test User",
			Email:    "Invalid",
			Password: "Password1",
		},
		expectedErrors: []expectedError{
			{"email", models.ValidationTypeInvalid},
		},
	},
	{
		description: "email required",
		model: &models.User{
			Name:     "Test User",
			Email:    "",
			Password: "Password1",
		},
		expectedErrors: []expectedError{
			{"email", models.ValidationTypeRequired},
			{"email", models.ValidationTypeInvalid},
		},
	},
}

func TestUserValidateCreate(t *testing.T) {
	userSetUp()
	defer userTearDown()

	createTestCases := []testCase{
		{
			description: "email not unique",
			model: &models.User{
				Name:     "Test User",
				Email:    existingUserEmail,
				Password: "Password1",
			},
			expectedErrors: []expectedError{
				{"email", models.ValidationTypeNotUnique},
			},
		},
		{
			description: "password required",
			model: &models.User{
				Name:     "Test User",
				Email:    "someone@example.com",
				Password: "",
			},
			expectedErrors: []expectedError{
				{"password", models.ValidationTypeRequired},
				{"password", models.ValidationTypeMinLength},
			},
		},
		{
			description: "password is too short",
			model: &models.User{
				Name:     "Test User",
				Email:    "someone@example.com",
				Password: "1234567",
			},
			expectedErrors: []expectedError{
				{"password", models.ValidationTypeMinLength},
			},
		},
	}

	var testCases []testCase
	testCases = append(testCases, commonTestCases...)
	testCases = append(testCases, createTestCases...)

	runValidationTest(t, testCases, func(user models.Model, id string) ([]models.ValidationError, error) {
		return user.ValidateCreate()
	})
}

func TestUserValidateUpdate(t *testing.T) {
	userSetUp()
	defer userTearDown()

	updateTestCases := []testCase{
		{
			description: "updating existing user",
			id:          existingUserID,
			model: &models.User{
				Name:     "Test User",
				Email:    existingUserEmail,
				Password: "Password1",
			},
			expectedErrors: []expectedError{},
		},
		{
			description: "password can be blank",
			id:          existingUserID,
			model: &models.User{
				Name:  "Test User",
				Email: existingUserEmail,
			},
			expectedErrors: []expectedError{},
		},
	}

	var testCases []testCase
	testCases = append(testCases, commonTestCases...)
	testCases = append(testCases, updateTestCases...)

	runValidationTest(t, testCases, func(user models.Model, id string) ([]models.ValidationError, error) {
		return user.ValidateUpdate(id)
	})
}
