package sqlboiler_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/sironfoot/go-twitter-bot/lib/sqlboiler"
)

func TestEntityGetByID(t *testing.T) {
	mustSetUp()
	defer mustTearDown()

	testUser, err := createTestUser()
	if err != nil {
		t.Fatal(err)
	}

	user := userEntity{}
	err = sqlboiler.EntityGetByID(&user, testUser.ID, db)
	if err != nil && err != sqlboiler.ErrEntityNotFound {
		t.Fatal(err)
	}

	if err == sqlboiler.ErrEntityNotFound {
		t.Errorf("Expected user record, but got ErrEntityNotFound")
	}

	if !usersAreSame(testUser, user) {
		t.Errorf("Expected user and actual user don't match,\nuser:\t\t%v,\ntestUser:\t%v", user, testUser)
	}

	// test mismatch ID types
	err = sqlboiler.EntityGetByID(&user, 123, db)
	if err == nil {
		t.Error("using mismatched ID types should return an error")
	}

	// record not found
	err = sqlboiler.EntityGetByID(&user, "089fe8c2-05ab-11e6-9e18-b32d264e490b", db)
	if err != sqlboiler.ErrEntityNotFound {
		t.Errorf("non-existant ID should return 'ErrEntityNotFound' error, error was: %s", err)
	}
}

func TestEntitySave(t *testing.T) {
	mustSetUp()
	defer mustTearDown()

	// CREATE record
	user := userEntity{
		Name:           "Test User",
		Email:          "test@example.com",
		HashedPassword: "Password1_",
		AuthToken:      sql.NullString{String: "my_auth_token", Valid: true},
		IsAdmin:        true,
		IsService:      false,
		DateCreated:    time.Now().UTC(),
	}

	err := sqlboiler.EntitySave(&user, db)
	if err != nil {
		t.Fatal(err)
	}

	if user.IsTransient() {
		t.Errorf("user is still in a transient state (non-saved), user ID: %s", user.ID)
	}

	// UPDATE record
	user.Name = "Updated User"
	user.Email = "updated@example.com"
	user.HashedPassword = "UpdatedPassword1"
	user.AuthToken = sql.NullString{String: "updated_auth_token", Valid: true}
	user.IsAdmin = false
	user.IsService = true
	user.DateCreated = time.Now().UTC()

	previousID := user.ID

	err = sqlboiler.EntitySave(&user, db)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != previousID {
		t.Errorf("user ID was updated from %s to %s", previousID, user.ID)
	}
}

func TestEntityDelete(t *testing.T) {
	mustSetUp()
	defer mustTearDown()

	// arrange - create the record
	userToSave := userEntity{
		Name:           "Test User",
		Email:          "test@example.com",
		HashedPassword: "Password1_",
		AuthToken:      sql.NullString{String: "my_auth_token", Valid: true},
		IsAdmin:        true,
		IsService:      false,
		DateCreated:    time.Now().UTC(),
	}

	err := sqlboiler.EntitySave(&userToSave, db)
	if err != nil {
		t.Fatal(err)
	}

	userID := userToSave.ID

	// get the user
	user := userEntity{}
	err = sqlboiler.EntityGetByID(&user, userID, db)
	if err != nil {
		t.Fatal(err)
	}

	// act - delete the user
	err = sqlboiler.EntityDelete(&user, db)
	if err != nil {
		t.Fatal(err)
	}

	// assert - check user exists
	err = sqlboiler.EntityGetByID(&userEntity{}, userID, db)
	if err != sqlboiler.ErrEntityNotFound {
		t.Errorf("should return error 'ErrEntityNotFound' for no record found, error was: %s", err)
	}
}
