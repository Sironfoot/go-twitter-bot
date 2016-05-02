package models

import (
	"strings"
	"time"
)

// Tweet represents a model for creating/updating a tweet posted to
// the create/update tweet REST API endpoints, complete with validation
type Tweet struct {
	Text     string    `json:"text"`
	PostOn   time.Time `json:"postOn"`
	IsPosted bool      `json:"isPosted"`
}

// Sanitise sanitises fields for the model, such as trimming whitespace
func (tweet *Tweet) Sanitise() {
	tweet.Text = strings.TrimSpace(tweet.Text)
}

// Validate provides validation logic for creating or updating a Tweet
func (tweet *Tweet) Validate() ([]ValidationError, error) {
	var validationErrors []ValidationError

	validationErrors = validateRequired(validationErrors, tweet.Text, "text")
	validationErrors = validateMaxLength(validationErrors, tweet.Text, 140, "text")

	return validationErrors, nil
}

// ValidateCreate provides validation logic for creating a new Tweet only
func (tweet *Tweet) ValidateCreate() ([]ValidationError, error) {
	validationErrors, err := tweet.Validate()
	if err != nil {
		return nil, err
	}

	return validationErrors, nil
}

// ValidateUpdate provides validation logic for updating an existing Tweet only,
// 'id' is the database primary key ID of the current Tweet being updated.
func (tweet *Tweet) ValidateUpdate(id string) ([]ValidationError, error) {
	validationErrors, err := tweet.Validate()
	if err != nil {
		return nil, err
	}

	return validationErrors, nil
}
