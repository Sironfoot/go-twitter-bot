package db_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/sironfoot/go-twitter-bot/data/db"
)

func usersAreSame(user1, user2 db.User) bool {
	return user1.ID == user2.ID &&
		user1.Name == user2.Name &&
		user1.Email == user2.Email &&
		user1.HashedPassword == user2.HashedPassword &&
		user1.AuthToken == user2.AuthToken &&
		user1.IsAdmin == user2.IsAdmin &&
		user1.IsService == user2.IsService &&
		user1.DateCreated.Unix() == user2.DateCreated.Unix() // accurate to nearest second
}

func TestUserFromID(t *testing.T) {
	mustSetUp()
	defer mustTearDown()

	// arrange
	testUser, err := createTestUser()

	// Non-existent record - Valid UUID
	// act
	user, err := db.UserFromID("ec24d9b2-fb39-11e5-9dcc-df8c5db12101")
	if err != nil && err != db.ErrEntityNotFound {
		t.Fatal(err)
	}

	// assert
	if err != db.ErrEntityNotFound {
		t.Errorf("user entity was returned from non-existent ID, userID: %s", user.ID)
	}

	// Non-existent record - Invalid UUID
	user, err = db.UserFromID("Nonsense")
	if err != nil && err != db.ErrEntityNotFound {
		t.Fatal(err)
	}

	// assert
	if err != db.ErrEntityNotFound {
		t.Errorf("user entity was returned from non-existent invalid ID, userID: %s", user.ID)
	}

	// Existing record
	// act
	user, err = db.UserFromID(testUser.ID)
	if err != nil && err != db.ErrEntityNotFound {
		t.Fatal(err)
	}

	// assert
	if err == db.ErrEntityNotFound {
		t.Errorf("Expected user record, but got ErrEntityNotFound")
	}

	if !usersAreSame(user, testUser) {
		t.Errorf("Expected user and actual user don't match,\nuser:\t\t%v,\ntestUser:\t%v", user, testUser)
	}
}

func TestUserFromEmail(t *testing.T) {
	mustSetUp()
	defer mustTearDown()

	// arrange
	testUser, err := createTestUser()

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

	if !usersAreSame(user, testUser) {
		t.Errorf("Expected user and actual user don't match")
	}
}

func TestUsersAll(t *testing.T) {
	mustSetUp()
	defer mustTearDown()

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
		_, err = createTestUserWithProperties(fmt.Sprintf("Test User %d", i), fmt.Sprintf("test_%d@example.com", i))
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

func TestUsersAllPaging(t *testing.T) {
	mustSetUp()
	defer mustTearDown()

	// setup
	max := 100
	limit := 20

	for i := 1; i <= max; i++ {
		_, err := createTestUserWithProperties(fmt.Sprintf("Test User %d", i), fmt.Sprintf("test_%d@example.com", i))
		if err != nil {
			t.Fatal(err)
		}
	}

	// arrange (1st page)
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
	if len(users) != limit {
		t.Errorf("expected %d records, actual was %d (1st page)", limit, len(users))
	}

	for i, user := range users {
		expectedEmail := fmt.Sprintf("test_%d@example.com", i+1)
		if user.Email != expectedEmail {
			t.Errorf("for user at index pos: %d (1st page), expected email %s actual email %s", i, expectedEmail, user.Email)
		}
	}

	// arrange (2nd page)
	query.After = users[len(users)-1].Email

	// act
	users, err = db.UsersAll(query)
	if err != nil {
		t.Fatal(err)
	}

	// assert
	if len(users) != limit {
		t.Errorf("expected %d records, actual was %d (2nd page)", limit, len(users))
	}

	for i, user := range users {
		expectedEmail := fmt.Sprintf("test_%d@example.com", (i+1)+limit)
		if user.Email != expectedEmail {
			t.Errorf("for user at index pos: %d (2nd page), expected email %s actual email %s", i, expectedEmail, user.Email)
		}
	}
}
