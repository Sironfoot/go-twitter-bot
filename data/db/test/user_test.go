package db_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/sironfoot/go-twitter-bot/data/db"
)

func TestUserFromID(t *testing.T) {
	if err := setUp(); err != nil {
		t.Fatal(err)
		return
	}

	defer func() {
		if err := tearDown(); err != nil {
			t.Fatal(err)
		}
	}()

	// arrange
	testUser, id, err := createTestUser()

	// Non-existent record - Valid UUID
	// act
	user, err := db.UserFromID("ec24d9b2-fb39-11e5-9dcc-df8c5db12101")
	if err != nil && err != db.ErrEntityNotFound {
		t.Fatal(err)
	}

	// assert
	if err != db.ErrEntityNotFound {
		t.Errorf("user entity was returned from non-existent ID, userID: %s", user.ID())
	}

	// Non-existent record - Invalid UUID
	user, err = db.UserFromID("Nonsense")
	if err != nil && err != db.ErrEntityNotFound {
		t.Fatal(err)
	}

	// assert
	if err != db.ErrEntityNotFound {
		t.Errorf("user entity was returned from non-existent invalid ID, userID: %s", user.ID())
	}

	// Existing record
	// act
	user, err = db.UserFromID(id)
	if err != nil && err != db.ErrEntityNotFound {
		t.Fatal(err)
	}

	// assert
	if err == db.ErrEntityNotFound {
		t.Errorf("Expected user record, but got ErrEntityNotFound")
	}

	if user.ID() != id ||
		user.Email != testUser.Email ||
		user.IsAdmin != testUser.IsAdmin ||
		user.DateCreated.Equal(testUser.DateCreated) {

		t.Errorf("Expected user and actual user don't match")
	}
}

func TestUserFromEmail(t *testing.T) {
	if err := setUp(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := tearDown(); err != nil {
			t.Fatal(err)
		}
	}()

	// arrange
	testUser, id, err := createTestUser()

	// Non-existent record
	// act
	user, err := db.UserFromEmail("Nonsense@example.com")
	if err != nil && err != db.ErrEntityNotFound {
		t.Fatal(err)
	}

	// assert
	if err != db.ErrEntityNotFound {
		t.Errorf("user entity was returned from non-existent email address, userEmail: %s", user.Email)
	}

	// Existing record
	// act
	user, err = db.UserFromEmail(testUser.Email)
	if err != nil && err != db.ErrEntityNotFound {
		t.Fatal(err)
	}

	// assert
	if err == db.ErrEntityNotFound {
		t.Errorf("Expected user record, but got ErrEntityNotFound")
	}

	if user.ID() != id ||
		user.Email != testUser.Email ||
		user.IsAdmin != testUser.IsAdmin ||
		user.DateCreated.Equal(testUser.DateCreated) {

		t.Errorf("Expected user and actual user don't match")
	}
}

func TestUsersAll(t *testing.T) {
	if err := setUp(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := tearDown(); err != nil {
			t.Fatal(err)
		}
	}()

	min := 20
	max := 100
	limit := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max-min) + min

	// test with empty database
	// arrange
	query := db.QueryAll{
		OrderBy:  db.UsersOrderByEmail,
		OrderAsc: true,
		Limit:    limit,
	}

	// act
	users, err := db.UsersAll(query)
	if err != nil {
		t.Fatal(err)
	}

	// assert
	if len(users) != 0 {
		t.Errorf("expected 0 records, actual was %d", len(users))
	}

	// multiple records
	// arrange
	for i := 1; i <= max; i++ {
		_, _, err = createTestUserWithEmail(fmt.Sprintf("test_%d@example.com", i))
		if err != nil {
			t.Fatal(err)
		}
	}

	// act
	users, err = db.UsersAll(query)
	if err != nil {
		t.Fatal(err)
	}

	// assert
	if len(users) != limit {
		t.Errorf("expected %d records, actual was %d", limit, len(users))
	}

	for i, user := range users {
		expectedEmail := fmt.Sprintf("test_%d@example.com", i+1)
		if user.Email != expectedEmail {
			t.Errorf("for user at index pos: %d, expected email %s actual email %s", i, expectedEmail, user.Email)
		}
	}
}
