package models

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

// Model is an interface for all model types
type Model interface {
	TrimFields()
	ValidateCreate() ([]ValidationError, error)
	ValidateUpdate(id string) ([]ValidationError, error)
}

// ValidationError contains information on a single validate error when validating a model
type ValidationError struct {
	FieldName string `json:"fieldName"`
	Type      string `json:"code"`
	Message   string `json:"message"`
}

const (
	// ValidationTypeRequired represents fields that are required
	ValidationTypeRequired = "required"

	// ValidationTypeInvalid represents fields that in an invalid
	// format, such as email addresses
	ValidationTypeInvalid = "invalid"

	// ValidationTypeMinLength represents fields that don't
	// meet the minimum required length of runes
	ValidationTypeMinLength = "min_length"

	// ValidationTypeMaxLength represents fields that exceed the
	// maximum allowed length of runes
	ValidationTypeMaxLength = "max_length"

	// ValidationTypeNotUnique represents fields that are required
	// to be unique in the database such as email addresses
	ValidationTypeNotUnique = "not_unique"
)

var isEmail = regexp.MustCompile(`(?i)^.+@.+\.[a-z]+$`)

func validateRequired(validationErrors []ValidationError, fieldValue, fieldName string) []ValidationError {
	if fieldValue == "" {
		validationErrors = append(validationErrors, ValidationError{
			FieldName: fieldName,
			Type:      ValidationTypeRequired,
			Message:   "'" + fieldName + "' field is required.",
		})
	}

	return validationErrors
}

func validateEmail(validationErrors []ValidationError, fieldValue, fieldName string) []ValidationError {
	if !isEmail.MatchString(fieldValue) {
		validationErrors = append(validationErrors, ValidationError{
			FieldName: fieldName,
			Type:      ValidationTypeInvalid,
			Message:   "'" + fieldName + "' is not a valid email address.",
		})
	}

	return validationErrors
}

func validateMaxLength(validationErrors []ValidationError, fieldValue string, maxLength int, fieldName string) []ValidationError {
	if utf8.RuneCountInString(fieldValue) > maxLength {
		validationErrors = append(validationErrors, ValidationError{
			FieldName: fieldName,
			Type:      ValidationTypeMaxLength,
			Message:   fmt.Sprintf("'%s' cannot be greater than %d characters.", fieldName, maxLength),
		})
	}

	return validationErrors
}

func validateMinLength(validationErrors []ValidationError, fieldValue string, minLength int, fieldName string) []ValidationError {
	if utf8.RuneCountInString(fieldValue) < minLength {
		validationErrors = append(validationErrors, ValidationError{
			FieldName: fieldName,
			Type:      ValidationTypeMinLength,
			Message:   fmt.Sprintf("'%s' cannot be less than %d characters.", fieldName, minLength),
		})
	}

	return validationErrors
}
