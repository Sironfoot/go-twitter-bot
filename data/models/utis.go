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

var isEmail = regexp.MustCompile(`(?i)^.+@.+\.[a-z]+$`)
