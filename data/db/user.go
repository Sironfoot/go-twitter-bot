package db

import "time"

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
		sql := "INSERT INTO users(email, hashed_password, is_admin, date_created) VALUES($1, $2, $3, $4) RETURNING id"

		statement, err := db.Prepare(sql)
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
		_, err := db.Exec("UPDATE users SET email = $2, hashed_password = $3 is_admin = $4, date_created = $5 WHERE id = $1",
			user.id, user.HashedPassword, user.Email, user.IsAdmin, user.DateCreated)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes the User from the database
func (user *User) Delete() error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", user.id)
	return err
}
