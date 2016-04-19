package models_test

import (
	"testing"
	"time"

	"github.com/sironfoot/go-twitter-bot/data/db"
	"github.com/sironfoot/go-twitter-bot/data/models"
	"github.com/sironfoot/go-twitter-bot/lib/sqlboiler"
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
				Email:          existingUserEmail,
				HashedPassword: "Password1",
				IsAdmin:        true,
				DateCreated:    time.Now(),
			}, nil
		}

		return db.User{}, sqlboiler.ErrEntityNotFound
	}
}

func userTearDown() {
	db.UserFromEmail = userFromEmail
}

var commonTestCases = []testCase{
	{
		description: "no errors",
		model: &models.User{
			Email:    "someone@example.com",
			Password: "Password1",
		},
		expectedErrors: map[string]string{},
	},
	{
		description: "email address invalid format",
		model: &models.User{
			Email:    "Invalid",
			Password: "Password1",
		},
		expectedErrors: map[string]string{"email": models.ValidationTypeInvalid},
	},
	{
		description: "email required",
		model: &models.User{
			Email:    "",
			Password: "Password1",
		},
		expectedErrors: map[string]string{"email": models.ValidationTypeRequired},
	},
}

func TestUserValidateCreate(t *testing.T) {
	userSetUp()
	defer userTearDown()

	createTestCases := []testCase{
		{
			description: "email not unique",
			model: &models.User{
				Email:    existingUserEmail,
				Password: "Password1",
			},
			expectedErrors: map[string]string{"email": models.ValidationTypeNotUnique},
		},
		{
			description: "password required",
			model: &models.User{
				Email:    "someone@example.com",
				Password: "",
			},
			expectedErrors: map[string]string{"password": models.ValidationTypeRequired},
		},
		{
			description: "password is too short",
			model: &models.User{
				Email:    "someone@example.com",
				Password: "1234567",
			},
			expectedErrors: map[string]string{"password": models.ValidationTypeMinLength},
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
				Email:    existingUserEmail,
				Password: "Password1",
			},
			expectedErrors: map[string]string{},
		},
		{
			description: "password can be blank",
			id:          existingUserID,
			model: &models.User{
				Email: existingUserEmail,
			},
			expectedErrors: map[string]string{},
		},
	}

	var testCases []testCase
	testCases = append(testCases, commonTestCases...)
	testCases = append(testCases, updateTestCases...)

	runValidationTest(t, testCases, func(user models.Model, id string) ([]models.ValidationError, error) {
		return user.ValidateUpdate(id)
	})
}
