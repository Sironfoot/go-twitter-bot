package db

import (
	"database/sql"
	"regexp"
	"time"
)

// User maps to users table
type User struct {
	id             string
	Email          string
	HashedPassword string
	IsAdmin        bool
	DateCreated    time.Time
}

// ID returns read-only Primary Key ID of User
func (user *User) ID() string {
	return user.id
}

// IsTransient determines if User record has been saved to the database,
// true means User struct has NOT been saved, false means it has.
func (user *User) IsTransient() bool {
	return len(user.id) == 0
}

// Save saves the User struct to the database.
func (user *User) Save() error {
	if user.IsTransient() {
		cmd := `INSERT INTO users(email, hashed_password, is_admin, date_created)
				VALUES($1, $2, $3, $4)
				RETURNING id`

		statement, err := db.Prepare(cmd)
		if err != nil {
			return err
		}
		defer statement.Close()

		err = statement.
			QueryRow(user.Email, user.HashedPassword, user.IsAdmin, user.DateCreated).
			Scan(&user.id)
		if err != nil {
			return err
		}
	} else {
		cmd := `UPDATE users
				SET email = $2, hashed_password = $3 is_admin = $4, date_created = $5
				WHERE id = $1`

		_, err := db.Exec(cmd, user.id, user.HashedPassword, user.Email, user.IsAdmin, user.DateCreated)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes the User from the database
func (user *User) Delete() error {
	cmd := `DELETE FROM users
			WHERE id = $1`

	_, err := db.Exec(cmd, user.id)
	return err
}

var isUUID = regexp.MustCompile(`(?i)^[a-f0-9]{8}\-[a-f0-9]{4}\-[a-f0-9]{4}\-[a-f0-9]{4}\-[a-f0-9]{12}$`)

// UserFromID returns a User record with given ID
func UserFromID(id string) (User, error) {
	var user User

	if !isUUID.MatchString(id) {
		return user, ErrEntityNotFound
	}

	cmd := `SELECT email, hashed_password, is_admin, date_created
			FROM users
			WHERE id = $1`

	err := db.QueryRow(cmd, id).
		Scan(&user.Email, &user.HashedPassword, &user.IsAdmin, &user.DateCreated)
	if err == sql.ErrNoRows {
		return user, ErrEntityNotFound
	} else if err != nil {
		return user, err
	}

	user.id = id
	return user, nil
}

// UserFromEmail returns the User record matching an email address
func UserFromEmail(email string) (User, error) {
	var user User

	cmd := `SELECT id, email, hashed_password, is_admin, date_created
			FROM users
			WHERE email = $1`

	err := db.QueryRow(cmd, email).
		Scan(&user.id, &user.Email, &user.HashedPassword, &user.IsAdmin, &user.DateCreated)
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
func UsersAll(query QueryAll) ([]User, error) {
	var users []User

	cmd := `SELECT id, email, hashed_password, is_admin, date_created
			FROM users
			ORDER BY $1
			LIMIT $2`

	orderBy := query.OrderBy
	if query.OrderAsc {
		orderBy += " ASC"
	} else {
		orderBy += " DESC"
	}

	rows, err := db.Query(cmd, orderBy, query.Limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.id, &user.Email, &user.HashedPassword, &user.IsAdmin, &user.DateCreated)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
