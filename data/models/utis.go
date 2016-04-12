package models

import "regexp"

// ValidationError contains information on a single validate error when validating a model
type ValidationError struct {
	FieldName string `json:"fieldName"`
	Message   string `json:"message"`
}

var isEmail = regexp.MustCompile(`(?i)^.+@.+\.[a-z]+$`)
