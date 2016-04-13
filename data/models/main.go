package models

import "regexp"

// Model is an interface for all model types
type Model interface {
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

	// ValidationTypeNotUnique represents fields that are required
	// to be unique in the database such as email addresses
	ValidationTypeNotUnique = "not_unique"
)

var isEmail = regexp.MustCompile(`(?i)^.+@.+\.[a-z]+$`)
