package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sironfoot/go-twitter-bot/data/db"
)

// User model returned by REST API
type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hasedPassword"`
	IsAdmin        bool      `json:"isAdmin"`
	DateCreated    time.Time `json:"dateCreated"`
}

// GetUsers = GET: /users
func GetUsers(res http.ResponseWriter, req *http.Request) {
	usersDB, err := db.UsersAll()
	if err != nil {
		panic(err)
	}

	var users []User
	for _, userDB := range usersDB {
		user := User{
			ID:             userDB.ID(),
			Email:          userDB.Email,
			HashedPassword: userDB.HashedPassword,
			IsAdmin:        userDB.IsAdmin,
			DateCreated:    userDB.DateCreated,
		}

		users = append(users, user)
	}

	data, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		panic(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}

// func GetUser(res http.ResponseWriter, req *http.Request) {
// 	vars := mux.Vars(req)
// 	userId := vars["userId"]
//
//     user, err := db.UserFromID(userId)
//     if err == db.ErrEntityNotFound {
//
//     }
// }
