package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

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

	router.UseC(func(next goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json; charset=utf-8")

			appContext := api.AppContext{}
			appContext.Settings = configuration
			ctx = context.WithValue(ctx, "appContext", &appContext)

			next.ServeHTTPC(ctx, res, req)
		})
	})

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
