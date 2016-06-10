package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/sironfoot/go-twitter-bot/data/api"
	"github.com/sironfoot/go-twitter-bot/data/db"
	"github.com/sironfoot/go-twitter-bot/lib/config"

	"goji.io"
	"goji.io/middleware"
	"goji.io/pat"
	"golang.org/x/net/context"
)

func main() {
	var configuration api.Config
	err := config.LoadWithCaching("config.json", "dev", &configuration)
	if err != nil {
		log.Fatal(err)
	}

	addr := flag.String("addr", configuration.AppSettings.ServerAddress, "Address to run server on")
	dbConn := flag.String("db", configuration.Database.ConnectionString, "Database connection string")
	flag.Parse()

	err = db.InitDB(*dbConn)
	if err != nil {
		log.Fatal(err)
	}

	router := goji.NewMux()

	// 1. Error handling
	router.UseC(func(next goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
			defer func() {
				r := recover()
				if r != nil {
					log.Println(r)

					response := api.MessageResponse{
						Message: "Internal Server Error",
					}

					res.Header().Set("Content-Type", "application/json; charset=utf-8")
					res.WriteHeader(http.StatusInternalServerError)

					data, jsonErr := json.MarshalIndent(response, "", "    ")
					if jsonErr == nil {
						res.Write(data)
					}
				}
			}()

			next.ServeHTTPC(ctx, res, req)
		})
	})

	// 2. Setup request context values
	router.UseC(func(next goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json; charset=utf-8")

			appContext := api.AppContext{}

			var configuration api.Config
			err := config.LoadWithCaching("config.json", "dev", &configuration)
			if err != nil {
				panic(err)
			}

			appContext.Settings = configuration
			ctx = context.WithValue(ctx, "appContext", &appContext)

			next.ServeHTTPC(ctx, res, req)
		})
	})

	// 3. Check authorisation/logged in state
	router.UseC(func(next goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
			defer func() {
				next.ServeHTTPC(ctx, res, req)
			}()

			accessToken := req.Header.Get("accessToken")
			if len(accessToken) == 0 {
				return
			}

			encryptedToken, err := base64.StdEncoding.DecodeString(accessToken)
			if err != nil {
				// not a valid Base64 string, force user to log in again
				return
			}

			appContext := ctx.Value("appContext").(*api.AppContext)
			block, err := aes.NewCipher([]byte(appContext.Settings.AppSettings.EncryptionKey))
			if err != nil {
				panic(err)
			}

			if len(encryptedToken) < aes.BlockSize {
				panic(fmt.Errorf("encryptedToken too short"))
			}

			iv := encryptedToken[:aes.BlockSize]
			encryptedToken = encryptedToken[aes.BlockSize:]

			cfb := cipher.NewCFBDecrypter(block, iv)
			cfb.XORKeyStream(encryptedToken, encryptedToken)

			token := string(encryptedToken)
			tokenParts := strings.Split(token, "_")
			if len(tokenParts) != 2 {
				// not a valid format (UserID_AuthToken), force user to log in again
				return
			}

			userID := tokenParts[0]
			authToken := tokenParts[1]

			user, err := db.UserFromID(userID)
			if err != nil {
				panic(err)
			}

			if user.AuthToken.Valid && user.AuthToken.String == authToken {
				appContext.AuthUser = &user
			}
		})
	})

	// 4. Check for 404
	var notFoundHandler = func(next goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
			if middleware.Handler(ctx) == nil {
				response := api.MessageResponse{
					Message: "Page Not Found",
				}

				res.WriteHeader(http.StatusNotFound)

				data, jsonErr := json.MarshalIndent(response, "", "    ")
				if jsonErr != nil {
					panic(jsonErr)
				}
				res.Write(data)
			} else {
				next.ServeHTTPC(ctx, res, req)
			}
		})
	}

	router.UseC(notFoundHandler)

	// 5. Serve response as JSON
	router.UseC(func(next goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
			next.ServeHTTPC(ctx, res, req)

			response := ctx.Value("appContext").(*api.AppContext).Response

			if response != nil {
				data, jsonErr := json.MarshalIndent(response, "", "    ")
				if jsonErr != nil {
					panic(jsonErr)
				}
				res.Write(data)
			}
		})
	})

	var mustBeLoggedIn = func(next goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
			appContext := ctx.Value("appContext").(*api.AppContext)

			if appContext.AuthUser == nil {
				res.WriteHeader(http.StatusUnauthorized)

				response := api.MessageResponse{
					Message: "Not authenticated. Please authenticate with PUT: /account/login",
				}
				data, jsonErr := json.MarshalIndent(response, "", "    ")
				if jsonErr != nil {
					panic(jsonErr)
				}
				res.Write(data)
			} else {
				next.ServeHTTPC(ctx, res, req)
			}
		})
	}

	router.HandleFuncC(pat.Get("/"), func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
		appContext := ctx.Value("appContext").(*api.AppContext)

		appContext.Response = api.MessageResponse{
			Message: "Hello from GoBot Data server",
		}
	})

	// Account
	account := goji.SubMux()
	account.UseC(notFoundHandler)
	router.HandleC(pat.New("/account/*"), account)

	account.HandleFuncC(pat.Put("/login"), api.AccountLogin)
	account.HandleFuncC(pat.Put("/logout"), api.AccountLogout)
	account.HandleFuncC(pat.Post("/signup"), api.AccountSignup)

	// Users
	users := goji.SubMux()
	users.UseC(notFoundHandler)
	users.UseC(mustBeLoggedIn)
	router.HandleC(pat.New("/users/*"), users)
	router.HandleC(pat.New("/users"), users)

	users.HandleFuncC(pat.Get(""), api.UsersAll)
	users.HandleFuncC(pat.Post(""), api.UserCreate)
	users.HandleFuncC(pat.Get("/:userID"), api.UserGet)
	users.HandleFuncC(pat.Put("/:userID"), api.UserUpdate)
	users.HandleFuncC(pat.Delete("/:userID"), api.UserDelete)

	// TwitterAccounts
	twitterAccounts := goji.SubMux()
	twitterAccounts.UseC(notFoundHandler)
	twitterAccounts.UseC(mustBeLoggedIn)
	router.HandleC(pat.New("/twitterAccounts/*"), twitterAccounts)
	router.HandleC(pat.New("/twitterAccounts"), twitterAccounts)

	twitterAccounts.HandleFuncC(pat.Get(""), api.TwitterAccountsAll)
	twitterAccounts.HandleFuncC(pat.Get("/:twitterAccountID"), api.TwitterAccountGet)

	twitterAccounts.HandleFuncC(pat.Get("/:twitterAccountID/tweets"), api.TwitterAccountGetWithTweets)
	twitterAccounts.HandleFuncC(pat.Post("/:twitterAccountID/tweets"), api.TwitterAccountTweetCreate)
	twitterAccounts.HandleFuncC(pat.Put("/:twitterAccountID/tweets/:tweetID"), api.TwitterAccountTweetUpdate)
	twitterAccounts.HandleFuncC(pat.Delete("/:twitterAccountID/tweets/:tweetID"), api.TwitterAccountTweetDelete)

	server := http.Server{
		Addr:    *addr,
		Handler: router,
	}

	log.Printf("GoBot Data Server running on %s...\n", *addr)
	log.Fatal(server.ListenAndServe())
}
