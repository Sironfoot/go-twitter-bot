package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/sironfoot/go-twitter-bot/data/api"
	"github.com/sironfoot/go-twitter-bot/data/db"
	"github.com/sironfoot/go-twitter-bot/lib/config"

	"goji.io"
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
			appContext := api.AppContext{}
			appContext.Settings = configuration
			ctx = context.WithValue(ctx, "appContext", &appContext)

			next.ServeHTTPC(ctx, res, req)
		})
	})

	router.UseC(func(next goji.Handler) goji.Handler {
		isRootPathMissingTrailingSlash := regexp.MustCompile(`(?i)^/[a-z0-9]+$`)

		return goji.HandlerFunc(func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
			if isRootPathMissingTrailingSlash.MatchString(req.URL.Path) {
				res.Header().Set("Location", req.URL.Path+"/")
				res.WriteHeader(http.StatusMovedPermanently)
			} else {
				next.ServeHTTPC(ctx, res, req)
			}
		})
	})

	router.UseC(func(next goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
			next.ServeHTTPC(ctx, res, req)

			response := ctx.Value("appContext").(*api.AppContext).Response

			res.Header().Set("Content-Type", "application/json")

			if response == nil {
				response = struct {
					Message string `json:"message"`
				}{"Response message was missing."}
				res.WriteHeader(http.StatusInternalServerError)
			}

			data, jsonErr := json.MarshalIndent(response, "", "    ")
			if jsonErr != nil {
				panic(jsonErr)
			}
			res.Write(data)
		})
	})

	router.HandleFuncC(pat.Get("/"), func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Hello from GoBot Data server\n")
	})

	// Account
	account := goji.SubMux()
	router.HandleC(pat.New("/account/*"), account)

	account.HandleFuncC(pat.Put("/login"), api.AccountLogin)
	account.HandleFuncC(pat.Put("/logout"), api.AccountLogout)
	account.HandleFuncC(pat.Post("/signup"), api.AccountSignup)

	// Users
	users := goji.SubMux()
	router.HandleC(pat.New("/users/*"), users)

	users.HandleFuncC(pat.Get("/"), api.UsersAll)
	users.HandleFuncC(pat.Post("/"), api.UserCreate)
	users.HandleFuncC(pat.Get("/:userID"), api.UserGet)
	users.HandleFuncC(pat.Put("/:userID"), api.UserUpdate)
	users.HandleFuncC(pat.Delete("/:userID"), api.UserDelete)

	// TwitterAccounts
	twitterAccounts := goji.SubMux()
	router.HandleC(pat.New("/twitterAccounts/*"), twitterAccounts)

	twitterAccounts.HandleFuncC(pat.Get("/"), api.TwitterAccountsAll)
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
