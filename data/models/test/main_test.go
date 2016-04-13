package models_test

import (
	"testing"

	"github.com/sironfoot/go-twitter-bot/data/models"
)

type testCase struct {
	description    string
	model          models.Model
	id             string
	expectedErrors map[string]string
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
				expectedType, ok := testCase.expectedErrors[validationError.FieldName]

				if !ok {
					t.Errorf("test case '%s': validationError field '%s' wasn't expected",
						testCase.description, validationError.FieldName)
				}

				if ok && validationError.Type != expectedType {
					t.Errorf("test case '%s': validationError field '%s' type (%s) doesn't match expected type (%s)",
						testCase.description, validationError.FieldName, validationError.Type, expectedType)
				}
			}
		}
	}
}
