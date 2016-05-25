package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"encoding/base64"

	"github.com/sironfoot/go-twitter-bot/data/db"
	"github.com/sironfoot/go-twitter-bot/data/models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

// AccountLogin = PUT: /account/login
func AccountLogin(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appContext := ctx.Value("appContext").(*AppContext)
	var login models.Login

	defer req.Body.Close()

	err := json.NewDecoder(io.LimitReader(req.Body, maxRequestLength)).Decode(&login)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		appContext.Response = MessageResponse{
			Message: "JSON request body was in invalid format",
		}
		return
	}

	login.Sanitise()
	validationErrors, err := login.Validate()
	if err != nil {
		panic(err)
	}

	if len(validationErrors) > 0 {
		res.WriteHeader(http.StatusBadRequest)
		appContext.Response = updateResponse{
			Message: "Login model is invalid.",
			Errors:  validationErrors,
		}
		return
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

		appContext.Response = updateResponse{
			Message: "Login model is invalid.",
			Errors:  validationErrors,
		}
		return
	} else if err != nil {
		panic(err)
	}

	// check password
	errHashCompare := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(login.Password))
	if errHashCompare != nil {
		if errHashCompare == bcrypt.ErrMismatchedHashAndPassword || errHashCompare == bcrypt.ErrHashTooShort {
			res.WriteHeader(http.StatusForbidden)

			var validationErrors []models.ValidationError
			validationErrors = append(validationErrors, models.ValidationError{
				FieldName: "password",
				Type:      models.ValidationTypeNotFound,
				Message:   fmt.Sprintf("Password not correct for user %s", login.Email),
			})

			appContext.Response = updateResponse{
				Message: "Login model is invalid.",
				Errors:  validationErrors,
			}
			return
		}

		panic(errHashCompare)
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

	err = user.Save()
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher([]byte(appContext.Settings.AppSettings.EncryptionKey))
	if err != nil {
		panic(err)
	}

	plaintextToken := []byte(user.ID + "_" + user.AuthToken.String)
	encryptedToken := make([]byte, aes.BlockSize+len(plaintextToken))

	// iv = initialization vector
	iv := encryptedToken[:aes.BlockSize]
	_, err = rand.Read(iv)
	if err != nil {
		panic(err)
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(encryptedToken[aes.BlockSize:], plaintextToken)
	accessToken := base64.StdEncoding.EncodeToString(encryptedToken)

	appContext.Response = struct {
		Message     string `json:"message"`
		AccessToken string `json:"accessToken"`
	}{
		Message:     ok,
		AccessToken: accessToken,
	}
}

// AccountLogout = PUT: /account/logout
func AccountLogout(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appContext := ctx.Value("appContext").(*AppContext)

	if appContext.AuthUser != nil {
		authUser := appContext.AuthUser

		authUser.AuthToken = sql.NullString{
			Valid: false,
		}

		err := authUser.Save()
		if err != nil {
			panic(err)
		}
	}

	appContext.Response = MessageResponse{
		Message: ok,
	}
}

// AccountSignup = POST/PUT: /account/signup
func AccountSignup(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotImplemented)
}
