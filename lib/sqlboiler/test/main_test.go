package sqlboiler_test

import (
	"database/sql"
	"io/ioutil"
	"time"

	_ "github.com/lib/pq" // initialise postgresql DB provider
	"github.com/sironfoot/go-twitter-bot/lib/sqlboiler"
)

// User maps to users table
type userEntity struct {
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
func (user *userEntity) IsTransient() bool {
	return len(user.ID) == 0
}

// MetaData returns meta data information about the User entity
func (user *userEntity) MetaData() sqlboiler.EntityMetaData {
	return sqlboiler.EntityMetaData{
		TableName:      "users",
		PrimaryKeyName: "id",
	}
}

var db *sql.DB

func mustSetUp() {
	tempDB, err := sql.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	_, err = tempDB.Exec("DROP DATABASE IF EXISTS go_twitter_bot_test;")
	if err != nil {
		panic(err)
	}

	_, err = tempDB.Exec("CREATE DATABASE go_twitter_bot_test;")
	if err != nil {
		panic(err)
	}

	tempDB.Close()

	db, err = sql.Open("postgres", "user=postgres dbname=go_twitter_bot_test sslmode=disable")
	if err != nil {
		panic(err)
	}

	createTables := "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";\n"

	file, err := ioutil.ReadFile("../../../data/sql/create-tables.sql")
	if err != nil {
		panic(err)
	}
	createTables += string(file)

	_, err = db.Exec(createTables)
	if err != nil {
		panic(err)
	}
}

func mustTearDown() {
	if err := db.Close(); err != nil {
		panic(err)
	}

	tempDB, err := sql.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	_, err = tempDB.Exec("DROP DATABASE IF EXISTS go_twitter_bot_test;")
	if err != nil {
		panic(err)
	}

	if err := tempDB.Close(); err != nil {
		panic(err)
	}

	return
}

func usersAreSame(user1, user2 userEntity) bool {
	return user1.ID == user2.ID &&
		user1.Name == user2.Name &&
		user1.Email == user2.Email &&
		user1.HashedPassword == user2.HashedPassword &&
		user1.AuthToken == user2.AuthToken &&
		user1.IsAdmin == user2.IsAdmin &&
		user1.IsService == user2.IsService &&
		user1.DateCreated.Unix() == user2.DateCreated.Unix() // accurate to nearest second
}

func createTestUser() (userEntity, error) {
	return createTestUserWithProperties("Test User", "test@example.com")
}

func createTestUserWithProperties(name, email string) (user userEntity, err error) {
	user = userEntity{
		Name:           name,
		Email:          email,
		HashedPassword: "Password1",
		AuthToken:      sql.NullString{},
		IsAdmin:        true,
		IsService:      false,
		DateCreated:    time.Now().UTC(),
	}

	createSQL := `INSERT INTO users(name, email, hashed_password, auth_token, is_admin, is_service, date_created)
                  VALUES($1, $2, $3, $4, $5, $6, $7)
                  RETURNING id`

	statement, err := db.Prepare(createSQL)
	if err != nil {
		return
	}
	defer statement.Close()

	err = statement.
		QueryRow(user.Name, user.Email, user.HashedPassword, user.AuthToken, user.IsAdmin, user.IsService, user.DateCreated).
		Scan(&user.ID)
	if err != nil {
		return
	}

	return
}
