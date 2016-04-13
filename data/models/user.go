package models

import (
	"strings"
	"unicode/utf8"

	"github.com/sironfoot/go-twitter-bot/data/db"
)

// User represents a model for creating/updating a user posted to
// the create/update user REST API endpoints, complete with validation
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

// Validate provides validation logic for creating or updating a User
func (user *User) Validate() ([]ValidationError, error) {
	var validationErrors []ValidationError

	if strings.TrimSpace(user.Email) == "" {
		validationErrors = append(validationErrors, ValidationError{
			FieldName: "email",
			Type:      ValidationTypeRequired,
			Message:   "'email' address is required.",
		})
	} else if !isEmail.MatchString(user.Email) {
		validationErrors = append(validationErrors, ValidationError{
			FieldName: "email",
			Type:      ValidationTypeInvalid,
			Message:   "'email' is not a valid email address.",
		})
	}

	return validationErrors, nil
}

// ValidateCreate provides validation logic for creating a new User only
func (user *User) ValidateCreate() ([]ValidationError, error) {
	validationErrors, err := user.Validate()
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(user.Password) == "" {
		validationErrors = append(validationErrors, ValidationError{
			FieldName: "password",
			Type:      ValidationTypeRequired,
			Message:   "'password' is required.",
		})
	} else if utf8.RuneCountInString(user.Password) < 8 {
		validationErrors = append(validationErrors, ValidationError{
			FieldName: "password",
			Type:      ValidationTypeMinLength,
			Message:   "'password' must be at least 8 characters.",
		})
	}

	_, err = db.UserFromEmail(user.Email)
	if err != nil && err != db.ErrEntityNotFound {
		return nil, err
	}

	if err != db.ErrEntityNotFound {
		validationErrors = append(validationErrors, ValidationError{
			FieldName: "email",
			Type:      ValidationTypeNotUnique,
			Message:   "'email' address is already in use.",
		})
	}

	return validationErrors, nil
}

// ValidateUpdate provides validation logic for updating an existing User only,
// 'id' is the database primary key ID of the current User being updated.
func (user *User) ValidateUpdate(id string) ([]ValidationError, error) {
	validationErrors, err := user.Validate()
	if err != nil {
		return nil, err
	}

	existingUser, err := db.UserFromEmail(user.Email)
	if err != db.ErrEntityNotFound {
		if err != nil {
			return nil, err
		}

		if existingUser.Email == user.Email && existingUser.ID != id {
			validationErrors = append(validationErrors, ValidationError{
				FieldName: "email",
				Type:      ValidationTypeNotUnique,
				Message:   "'email' address is already in use.",
			})
		}
	}

	return validationErrors, nil
}
