package models

import "testing"

func TestUserValidateCreate(t *testing.T) {
	setUp()
	defer tearDown()

	testCases := []testCase{
		// 1. no errors
		{
			model: &User{
				Email:    "someone@example.com",
				Password: "Password1",
			},
			expectedErrors: map[string]string{},
		},

		// 2. email not unique
		{
			model: &User{
				Email:    "existing@example.com",
				Password: "Password1",
			},
			expectedErrors: map[string]string{"email": "not_unique"},
		},

		// 3. email address invalid format
		{
			model: &User{
				Email:    "Invalid",
				Password: "Password1",
			},
			expectedErrors: map[string]string{"email": "invalid"},
		},

		// 4. email and password are blank
		{
			model: &User{
				Email:    "",
				Password: "",
			},
			expectedErrors: map[string]string{"email": "required", "password": "required"},
		},

		// 5. password is too short (min 8 chars)
		{
			model: &User{
				Email:    "someone@example.com",
				Password: "abcd",
			},
			expectedErrors: map[string]string{"password": "min_length"},
		},
	}

	runValidationTest(t, testCases, func(user Model, id string) ([]ValidationError, error) {
		return user.ValidateCreate()
	})
}
