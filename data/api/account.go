package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"crypto/aes"
	"crypto/rand"

	"encoding/base64"

	"github.com/sironfoot/go-twitter-bot/data/db"
	"github.com/sironfoot/go-twitter-bot/data/models"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AccountLogin = POST/PUT: /account/login
func AccountLogin(res http.ResponseWriter, req *http.Request) interface{} {
	var login models.Login

	err := json.NewDecoder(req.Body).Decode(&login)
	if err != nil {
		panic(err)
	}
	req.Body.Close()

	login.Sanitise()
	validationErrors, err := login.Validate()
	if err != nil {
		panic(err)
	}

	if len(validationErrors) > 0 {
		res.WriteHeader(http.StatusBadRequest)
		return updateResponse{
			Message: "Login model is invalid.",
			Errors:  validationErrors,
		}
	}

	// check user exists for email
	user, err := db.UserFromEmail(login.Email)
	if err == db.ErrEntityNotFound {
		res.WriteHeader(http.StatusForbidden)

		var validationErrors []models.ValidationError
		validationErrors = append(validationErrors, models.ValidationError{
			FieldName: "email",
			Type:      models.ValidationTypeNotFound,
			Message:   fmt.Sprintf("User not found on email %s", login.Email),
		})

		return updateResponse{
			Message: "Login model is invalid.",
			Errors:  validationErrors,
		}
	} else if err != nil {
		panic(err)
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(login.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword || err == bcrypt.ErrHashTooShort {
			res.WriteHeader(http.StatusForbidden)

			var validationErrors []models.ValidationError
			validationErrors = append(validationErrors, models.ValidationError{
				FieldName: "password",
				Type:      models.ValidationTypeNotFound,
				Message:   fmt.Sprintf("Password not correct for user %s", login.Email),
			})

			return updateResponse{
				Message: "Login model is invalid.",
				Errors:  validationErrors,
			}
		}

		panic(err)
	}

	randomBytes := make([]byte, 32)
	_, err = rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	user.AuthToken = sql.NullString{
		String: base64.StdEncoding.EncodeToString(randomBytes),
		Valid:  true,
	}

	cypher, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		panic(err)
	}

	plaintextTokenBytes := []byte(user.ID + "_" + user.AuthToken.String)
	var encryptedTokenBytes []byte

	cypher.Encrypt(encryptedTokenBytes, plaintextTokenBytes)
	accessToken := base64.StdEncoding.EncodeToString(encryptedTokenBytes)

	return struct {
		Message     string `json:"message"`
		AccessToken string `json:"accessToken"`
	}{
		Message:     "",
		AccessToken: accessToken,
	}
}

// AccountLogout = POST/PUT: /account/logout
func AccountLogout(res http.ResponseWriter, req *http.Request) interface{} {
	return ""
}

// AccountSignup = POST/PUT: /account/signup
func AccountSignup(res http.ResponseWriter, req *http.Request) interface{} {
	return ""
}
