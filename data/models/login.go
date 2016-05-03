package models

import "strings"

// Login represents a model for logggin into the system
// using REST API endpoints, complete with validation
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Sanitise sanitises fields for the model, such as trimming whitespace
func (login *Login) Sanitise() {
	login.Email = strings.TrimSpace(login.Email)
}

// Validate provides validation logic for creating a new User only
func (login *Login) Validate() ([]ValidationError, error) {
	var validationErrors []ValidationError
	validationErrors = validateRequired(validationErrors, login.Email, "email")
	validationErrors = validateRequired(validationErrors, login.Password, "password")

	return validationErrors, nil
}
