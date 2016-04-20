package models_test

import (
	"testing"

	"github.com/sironfoot/go-twitter-bot/data/models"
)

type testCase struct {
	description    string
	model          models.Model
	id             string
	expectedErrors []expectedError
}

type expectedError struct {
	fieldName string
	typeName  string
}

func runValidationTest(t *testing.T, testCases []testCase, runValidation func(model models.Model, id string) ([]models.ValidationError, error)) {
	for _, testCase := range testCases {
		validationErrors, err := runValidation(testCase.model, testCase.id)
		if err != nil {
			t.Fatal(err)
		}

		if len(validationErrors) != len(testCase.expectedErrors) {
			t.Errorf("test case '%s': expected %d validation error(s) but got %d: %s",
				testCase.description, len(testCase.expectedErrors), len(validationErrors), validationErrors)
		} else {
			for _, validationError := range validationErrors {
				found := false
				for _, expectedError := range testCase.expectedErrors {
					if expectedError.fieldName == validationError.FieldName &&
						expectedError.typeName == validationError.Type {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("test case '%s': validationError field '%s(%s)' wasn't found",
						testCase.description, validationError.FieldName, validationError.Type)
				}
			}
		}
	}
}
