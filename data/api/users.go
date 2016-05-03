package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sironfoot/go-twitter-bot/data/db"
	"github.com/sironfoot/go-twitter-bot/data/models"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	IsAdmin     bool      `json:"isAdmin"`
	IsService   bool      `json:"isService"`
	DateCreated time.Time `json:"dateCreated"`
}

// UsersAll = GET: /users
func UsersAll(res http.ResponseWriter, req *http.Request) interface{} {
	defaults := getPagingDefaults(db.UsersOrderByDateCreated, false, db.UsersSortableColumns)
	paging, err := ExtractAndValidatePagingInfo(req, defaults)
	if err != nil {
		return messageResponse{err.Error()}
	}

	model := struct {
		pagedResponse
		Users []user `json:"users"`
	}{}

	usersDB, totalRecords, err := db.UsersAll(paging)
	if err != nil {
		panic(err)
	}

	var users []user
	for _, userDB := range usersDB {
		users = append(users, user{
			ID:          userDB.ID,
			Name:        userDB.Name,
			Email:       userDB.Email,
			IsAdmin:     userDB.IsAdmin,
			IsService:   userDB.IsService,
			DateCreated: userDB.DateCreated,
		})
	}

	model.Message = ok
	model.Page = paging.Page
	model.RecordsPerPage = paging.RecordsPerPage
	model.TotalRecords = totalRecords
	model.Users = users

	return model
}

// UserGet = GET: /users/[userID]
func UserGet(res http.ResponseWriter, req *http.Request) interface{} {
	vars := mux.Vars(req)
	userID := vars["userID"]

	userDB, err := db.UserFromID(userID)
	if err == db.ErrEntityNotFound {
		return messageResponse{
			Message: fmt.Sprintf("User not found on ID: %s", userID),
		}
	} else if err != nil {
		panic(err)
	}

	model := struct {
		messageResponse
		User user `json:"user"`
	}{}

	model.Message = ok
	model.User = user{
		ID:          userDB.ID,
		Name:        userDB.Name,
		Email:       userDB.Email,
		IsAdmin:     userDB.IsAdmin,
		IsService:   userDB.IsService,
		DateCreated: userDB.DateCreated,
	}

	return model
}

// UserCreate = POST: /users
func UserCreate(res http.ResponseWriter, req *http.Request) interface{} {
	var newUser models.User

	err := json.NewDecoder(req.Body).Decode(&newUser)
	if err != nil {
		panic(err)
	}
	req.Body.Close()

	newUser.Sanitise()
	validationErrors, err := newUser.ValidateCreate()
	if err != nil {
		panic(err)
	}

	model := createResponse{}

	if len(validationErrors) > 0 {
		model.Message = "User model is invalid."
		model.Errors = validationErrors
		res.WriteHeader(http.StatusBadRequest)
		return model
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)
	if err != nil {
		panic(err)
	}

	user := db.User{
		Name:           newUser.Name,
		Email:          newUser.Email,
		HashedPassword: string(hashedPassword),
		DateCreated:    time.Now().UTC(),
		IsAdmin:        newUser.IsAdmin,
		IsService:      newUser.IsService,
	}

	err = user.Save()
	if err != nil {
		panic(err)
	}

	model.Message = ok
	model.ID = &user.ID
	res.WriteHeader(http.StatusCreated)

	return model
}

// UserUpdate = PUT: /users/{userID}
func UserUpdate(res http.ResponseWriter, req *http.Request) interface{} {
	var updateUser models.User

	err := json.NewDecoder(req.Body).Decode(&updateUser)
	if err != nil {
		panic(err)
	}
	req.Body.Close()

	vars := mux.Vars(req)
	userID := vars["userID"]

	user, err := db.UserFromID(userID)
	if err != nil && err != db.ErrEntityNotFound {
		panic(err)
	}

	if err == db.ErrEntityNotFound {
		res.WriteHeader(http.StatusNotFound)
		return messageResponse{
			Message: fmt.Sprintf("User not found on ID: %s", userID),
		}
	}

	updateUser.Sanitise()
	validationErrors, err := updateUser.ValidateUpdate(userID)
	if err != nil {
		panic(err)
	}

	if len(validationErrors) > 0 {
		res.WriteHeader(http.StatusBadRequest)
		return updateResponse{
			Message: "User model is invalid.",
			Errors:  validationErrors,
		}
	}

	user.Email = updateUser.Email
	user.IsAdmin = updateUser.IsAdmin

	if updateUser.Password != "" {
		hashedPassword, bcryptErr := bcrypt.GenerateFromPassword([]byte(updateUser.Password), 12)
		if bcryptErr != nil {
			panic(bcryptErr)
		}

		user.HashedPassword = string(hashedPassword)
	}

	err = user.Save()
	if err != nil {
		panic(err)
	}

	return messageResponse{
		Message: ok,
	}
}

// UserDelete = DELETE: /users/{userID}
func UserDelete(res http.ResponseWriter, req *http.Request) interface{} {
	vars := mux.Vars(req)
	userID := vars["userID"]

	user, err := db.UserFromID(userID)
	if err != nil && err != db.ErrEntityNotFound {
		panic(err)
	}

	if err == db.ErrEntityNotFound {
		res.WriteHeader(http.StatusNotFound)
		return messageResponse{
			Message: fmt.Sprintf("User not found on ID: %s", userID),
		}
	}

	err = user.Delete()
	if err != nil {
		panic(err)
	}

	return messageResponse{
		Message: ok,
	}
}
