package models

import (
	"strings"

	"github.com/sironfoot/go-twitter-bot/data/db"
)

// User represents a model for creating/updating a user posted to
// the create/update user REST API endpoints, complete with validation
type User struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsAdmin   bool   `json:"isAdmin"`
	IsService bool   `json:"isService"`
}

// TrimFields trims whitespace from start and end of fields
// that are appropriate for trimming
func (user *User) TrimFields() {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
}

// Validate provides validation logic for creating or updating a User
func (user *User) Validate() ([]ValidationError, error) {
	var validationErrors []ValidationError

	validationErrors = validateRequired(validationErrors, user.Name, "name")
	validationErrors = validateMaxLength(validationErrors, user.Name, 50, "name")

	validationErrors = validateRequired(validationErrors, user.Email, "email")
	validationErrors = validateMaxLength(validationErrors, user.Email, 200, "email")
	validationErrors = validateEmail(validationErrors, user.Email, "email")

	return validationErrors, nil
}

// ValidateCreate provides validation logic for creating a new User only
func (user *User) ValidateCreate() ([]ValidationError, error) {
	validationErrors, err := user.Validate()
	if err != nil {
		return nil, err
	}

	validationErrors = validateRequired(validationErrors, user.Password, "password")
	validationErrors = validateMinLength(validationErrors, user.Password, 8, "password")

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
