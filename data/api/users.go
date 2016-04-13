package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sironfoot/go-twitter-bot/data/db"
	"github.com/sironfoot/go-twitter-bot/data/models"
)

const ok = "OK"

// User model returned by REST API
type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	IsAdmin     bool      `json:"isAdmin"`
	DateCreated time.Time `json:"dateCreated"`
}

// UsersAll = GET: /users
func UsersAll(res http.ResponseWriter, req *http.Request) {
	usersDB, err := db.UsersAll(db.QueryAll{})
	if err != nil {
		panic(err)
	}

	var users []User
	for _, userDB := range usersDB {
		user := User{
			ID:          userDB.ID,
			Email:       userDB.Email,
			IsAdmin:     userDB.IsAdmin,
			DateCreated: userDB.DateCreated,
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

// UserGet = GET: /users/[userID]
func UserGet(res http.ResponseWriter, req *http.Request) {
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

	model.Message = ok
	model.User = User{
		ID:          user.ID,
		Email:       user.Email,
		IsAdmin:     user.IsAdmin,
		DateCreated: user.DateCreated,
	}
}

// UserCreate = POST: /users
func UserCreate(res http.ResponseWriter, req *http.Request) {
	var newUser models.User

	err := json.NewDecoder(req.Body).Decode(&newUser)
	if err != nil {
		panic(err)
	}
	req.Body.Close()

	res.Header().Set("Content-Type", "application/json")

	validationErrors, err := newUser.ValidateCreate()
	if err != nil {
		panic(err)
	}

	model := struct {
		Message string                   `json:"message"`
		Errors  []models.ValidationError `json:"errors"`
	}{}

	if len(validationErrors) > 0 {
		model.Message = "User model is invalid."
		model.Errors = validationErrors
		res.WriteHeader(http.StatusBadRequest)
	} else {
		model.Message = ok
		res.WriteHeader(http.StatusCreated)
	}

	data, err := json.MarshalIndent(model, "", "   ")
	if err != nil {
		panic(err)
	}

	res.Write(data)
}

// UserUpdate = PUT: /users/{userID}
func UserUpdate(res http.ResponseWriter, req *http.Request) {
	var updateUser models.User

	err := json.NewDecoder(req.Body).Decode(&updateUser)
	if err != nil {
		panic(err)
	}
	req.Body.Close()

	vars := mux.Vars(req)
	userID := vars["userID"]

	res.Header().Set("Content-Type", "application/json")

	model := struct {
		Message string                   `json:"message"`
		Errors  []models.ValidationError `json:"errors"`
	}{}

	defer func() {
		data, jsonErr := json.MarshalIndent(model, "", "   ")
		if jsonErr != nil {
			panic(jsonErr)
		}

		res.Write(data)
	}()

	user, err := db.UserFromID(userID)
	if err != nil && err != db.ErrEntityNotFound {
		panic(err)
	}

	if err == db.ErrEntityNotFound {
		model.Message = fmt.Sprintf("User not found on ID: %s", userID)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	validationErrors, err := updateUser.ValidateUpdate(userID)

	if len(validationErrors) > 0 {
		model.Message = "User model is invalid."
		model.Errors = validationErrors
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	user.Email = updateUser.Email
	user.HashedPassword = updateUser.Password
	user.IsAdmin = updateUser.IsAdmin

	err = user.Save()
	if err != nil {
		panic(err)
	}

	model.Message = ok
}

// UserDelete = DELETE: /users/{userID}
func UserDelete(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := vars["userID"]

	model := struct {
		Message string `json:"message"`
	}{}

	defer func() {
		data, jsonErr := json.MarshalIndent(model, "", "   ")
		if jsonErr != nil {
			panic(jsonErr)
		}

		res.Write(data)
	}()

	user, err := db.UserFromID(userID)
	if err != nil && err != db.ErrEntityNotFound {
		panic(err)
	}

	if err == db.ErrEntityNotFound {
		model.Message = fmt.Sprintf("User not found on ID: %s", userID)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	err = user.Delete()
	if err != nil {
		panic(err)
	}

	model.Message = ok
}
