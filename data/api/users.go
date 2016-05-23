package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"goji.io/pat"

	"github.com/sironfoot/go-twitter-bot/data/db"
	"github.com/sironfoot/go-twitter-bot/data/models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
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
func UsersAll(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appContext := ctx.Value("appContext").(*AppContext)

	defaults := getPagingDefaults(db.UsersOrderByDateCreated, false, db.UsersSortableColumns)
	paging, err := ExtractAndValidatePagingInfo(req, defaults)
	if err != nil {
		appContext.Response = MessageResponse{err.Error()}
		return
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

	appContext.Response = model
}

// UserGet = GET: /users/:userID
func UserGet(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appContext := ctx.Value("appContext").(*AppContext)
	userID := pat.Param(ctx, "userID")

	userDB, err := db.UserFromID(userID)
	if err == db.ErrEntityNotFound {
		appContext.Response = MessageResponse{
			Message: fmt.Sprintf("User not found on ID: %s", userID),
		}
		return
	} else if err != nil {
		panic(err)
	}

	model := struct {
		MessageResponse
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

	appContext.Response = model
}

// UserCreate = POST: /users
func UserCreate(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appContext := ctx.Value("appContext").(*AppContext)
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
		appContext.Response = model

		res.WriteHeader(http.StatusBadRequest)
		return
	}

	bcryptWorkFactor := appContext.Settings.AppSettings.BCryptWorkFactor
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcryptWorkFactor)
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

	appContext.Response = model
}

// UserUpdate = PUT: /users/:userID
func UserUpdate(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appContext := ctx.Value("appContext").(*AppContext)
	var updateUser models.User

	err := json.NewDecoder(req.Body).Decode(&updateUser)
	if err != nil {
		panic(err)
	}
	req.Body.Close()

	userID := pat.Param(ctx, "userID")

	user, err := db.UserFromID(userID)
	if err != nil && err != db.ErrEntityNotFound {
		panic(err)
	}

	if err == db.ErrEntityNotFound {
		res.WriteHeader(http.StatusNotFound)
		appContext.Response = MessageResponse{
			Message: fmt.Sprintf("User not found on ID: %s", userID),
		}
		return
	}

	updateUser.Sanitise()
	validationErrors, err := updateUser.ValidateUpdate(userID)
	if err != nil {
		panic(err)
	}

	if len(validationErrors) > 0 {
		res.WriteHeader(http.StatusBadRequest)
		appContext.Response = updateResponse{
			Message: "User model is invalid.",
			Errors:  validationErrors,
		}
		return
	}

	user.Email = updateUser.Email
	user.IsAdmin = updateUser.IsAdmin

	if updateUser.Password != "" {
		bcryptWorkFactor := appContext.Settings.AppSettings.BCryptWorkFactor
		hashedPassword, bcryptErr := bcrypt.GenerateFromPassword([]byte(updateUser.Password), bcryptWorkFactor)
		if bcryptErr != nil {
			panic(bcryptErr)
		}

		user.HashedPassword = string(hashedPassword)
	}

	err = user.Save()
	if err != nil {
		panic(err)
	}

	appContext.Response = MessageResponse{
		Message: ok,
	}
}

// UserDelete = DELETE: /users/:userID
func UserDelete(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appContext := ctx.Value("appContext").(*AppContext)
	userID := pat.Param(ctx, "userID")

	user, err := db.UserFromID(userID)
	if err != nil && err != db.ErrEntityNotFound {
		panic(err)
	}

	if err == db.ErrEntityNotFound {
		res.WriteHeader(http.StatusNotFound)
		appContext.Response = MessageResponse{
			Message: fmt.Sprintf("User not found on ID: %s", userID),
		}
		return
	}

	err = user.Delete()
	if err != nil {
		panic(err)
	}

	appContext.Response = MessageResponse{
		Message: ok,
	}
}
