package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/sironfoot/go-twitter-bot/data/api"
	"github.com/sironfoot/go-twitter-bot/data/db"
	"github.com/sironfoot/go-twitter-bot/lib/config"

	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
)

// Config represents a configuration settings for the app
type Config struct {
	Database    Database    `json:"database"`
	AppSettings AppSettings `json:"appSettings"`
}

// Database represents database configuration settings for the app
type Database struct {
	DriverName       string `json:"driverName"`
	ConnectionString string `json:"connectionString"`
}

// AppSettings represents general application settings for the app
type AppSettings struct {
	ServerAddress string `json:"serverAddress"`
	EncryptionKey string `json:"encryptionKey"`
}

func main() {
	var configuration Config
	err := config.Load("config.json", "dev", &configuration)
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

	// router.UseC(func(inner goji.Handler) goji.Handler {
	// 	return goji.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// 		log.Print("A: before")
	// 		inner.ServeHTTPC(ctx, w, r)
	// 		log.Print("A: after")
	// 	})
	// })

	router.HandleFuncC(pat.Get("/"), func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Hello from GoBot Data server\n")
	})

	// Account
	account := goji.SubMux()
	router.HandleC(pat.New("/account/*"), account)

	account.HandleFuncC(pat.Put("/login"), wrapJSON(api.AccountLogin))
	account.HandleFuncC(pat.Put("/logout"), wrapJSON(api.AccountLogout))
	account.HandleFuncC(pat.Post("/signup"), wrapJSON(api.AccountSignup))

	// Users
	users := goji.SubMux()
	router.HandleC(pat.New("/users/*"), users)

	users.HandleFuncC(pat.Get("/"), wrapJSON(api.UsersAll))
	users.HandleFuncC(pat.Post("/"), wrapJSON(api.UserCreate))
	users.HandleFuncC(pat.Get("/:userID"), wrapJSON(api.UserGet))
	users.HandleFuncC(pat.Put("/:userID"), wrapJSON(api.UserUpdate))
	users.HandleFuncC(pat.Delete(":/userID"), wrapJSON(api.UserDelete))

	// TwitterAccounts
	twitterAccounts := goji.SubMux()
	router.HandleC(pat.New("/twitterAccounts/*"), twitterAccounts)

	twitterAccounts.HandleFuncC(pat.Get("/"), wrapJSON(api.TwitterAccountsAll))
	twitterAccounts.HandleFuncC(pat.Get("/:twitterAccountID"), wrapJSON(api.TwitterAccountGet))

	twitterAccounts.HandleFuncC(pat.Get("/:twitterAccountID/tweets"), wrapJSON(api.TwitterAccountGetWithTweets))
	twitterAccounts.HandleFuncC(pat.Post("/:twitterAccountID/tweets"), wrapJSON(api.TwitterAccountTweetCreate))
	twitterAccounts.HandleFuncC(pat.Put("/:twitterAccountID/tweets/:tweetID"), wrapJSON(api.TwitterAccountTweetUpdate))
	twitterAccounts.HandleFuncC(pat.Delete("/:twitterAccountID/tweets/:tweetID"), wrapJSON(api.TwitterAccountTweetDelete))

	server := http.Server{
		Addr:    *addr,
		Handler: router,
	}

	log.Printf("GoBot Data Server running on %s...\n", *addr)
	log.Fatal(server.ListenAndServe())
}

func wrapJSON(apiFunc func(context.Context, http.ResponseWriter, *http.Request) interface{}) goji.HandlerFunc {
	return func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
		response := apiFunc(ctx, res, req)

		defer func() {
			res.Header().Set("Content-Type", "application/json")

			data, jsonErr := json.MarshalIndent(response, "", "    ")
			if jsonErr != nil {
				panic(jsonErr)
			}
			res.Write(data)
		}()
	}
}
