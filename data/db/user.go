package db

import (
	"database/sql"
	"time"

	"github.com/sironfoot/go-twitter-bot/lib/sqlboiler"
)

// User maps to users table
type User struct {
	ID             string         `db:"id"`
	Name           string         `db:"name"`
	Email          string         `db:"email"`
	HashedPassword string         `db:"hashed_password"`
	AuthToken      sql.NullString `db:"auth_token"`
	IsAdmin        bool           `db:"is_admin"`
	IsService      bool           `db:"is_service"`
	DateCreated    time.Time      `db:"date_created"`
}

// IsTransient determines if User record has been saved to the database,
// true means User struct has NOT been saved, false means it has.
func (user *User) IsTransient() bool {
	return len(user.ID) == 0
}

// MetaData returns meta data information about the User entity
func (user *User) MetaData() sqlboiler.EntityMetaData {
	return sqlboiler.EntityMetaData{
		TableName:      "users",
		PrimaryKeyName: "id",
	}
}

// UserSave saves the User struct to the database.
var UserSave = func(user *User) error {
	return sqlboiler.EntitySave(user, db)
}

// Save saves the User struct to the database.
func (user *User) Save() error {
	return UserSave(user)
}

// UserDelete deletes the User from the database
var UserDelete = func(user *User) error {
	return sqlboiler.EntityDelete(user, db)
}

// Delete deletes the User from the database
func (user *User) Delete() error {
	return UserDelete(user)
}

// UserFromID returns a User record with given ID
var UserFromID = func(id string) (User, error) {
	var user User

	if !isUUID.MatchString(id) {
		return user, ErrEntityNotFound
	}

	err := sqlboiler.EntityGetByID(&user, id, db)
	if err == sqlboiler.ErrEntityNotFound {
		return user, ErrEntityNotFound
	}
	return user, err
}

// UserFromEmail returns the User record matching an email address
var UserFromEmail = func(email string) (User, error) {
	var user User

	cmd := `SELECT id, ` + sqlboiler.GetColumnListString(&user, "") + `
			FROM users
			WHERE email = $1`

	err := dbx.QueryRowx(cmd, email).StructScan(&user)
	if err == sql.ErrNoRows {
		return user, ErrEntityNotFound
	} else if err != nil {
		return user, err
	}

	return user, nil
}

const (
	// UsersOrderByDateCreated is for ordering users by DateCreated
	UsersOrderByDateCreated = "date_created"
	// UsersOrderByEmail is for ordering users by Email address
	UsersOrderByEmail = "email"
)

// UsersAll returns all User records from the database
var UsersAll = func(query PagingInfo) ([]User, int, error) {
	var users []User
	recordCount := 0

	cmd := "SELECT id, " + sqlboiler.GetColumnListString(&User{}, "") + " " +
		"FROM users " +
		"ORDER BY $1 " +
		"LIMIT $2 OFFSET $3"

	rows, err := dbx.Queryx(cmd, query.BuildOrderBy(), query.Limit(), query.Offset())
	if err != nil {
		return nil, recordCount, err
	}

	defer rows.Close()

	for rows.Next() {
		user := User{}
		err = rows.StructScan(&user)
		if err != nil {
			return nil, recordCount, err
		}

		users = append(users, user)
	}

	err = dbx.Get(&recordCount, "SELECT COUNT(*) FROM users")
	return users, recordCount, err
}
