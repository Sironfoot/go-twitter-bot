package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

// GetUser = GET: /users/[userID]
func GetUser(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := vars["userID"]

	model := struct {
		Message string `json:"message"`
		User    User   `json:"user"`
	}{}

	res.Header().Set("Content-Type", "application/json")

	defer func() {
		data, err := json.MarshalIndent(model, "", "   ")
		if err != nil {
			panic(err)
		}

		res.Write(data)
	}()

	user, err := db.UserFromID(userID)
	if err == db.ErrEntityNotFound {
		model.Message = fmt.Sprintf("User not found on ID: %s", userID)
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		panic(err)
	}

	model.Message = "OK"
	model.User = User{
		ID:             user.ID(),
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		IsAdmin:        user.IsAdmin,
		DateCreated:    user.DateCreated,
	}
}
