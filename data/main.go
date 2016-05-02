package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sironfoot/go-twitter-bot/data/api"
	"github.com/sironfoot/go-twitter-bot/data/db"
)

var addr = flag.String("addr", "localhost:7001", "Address to run server on")

func main() {
	flag.Parse()

	err := db.InitDB("user=postgres dbname=go_twitter_bot sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Hello from GoBot Data server\n")
	})

	router.HandleFunc("/users", wrapJSON(api.UsersAll)).
		Methods("GET")
	router.HandleFunc("/users/{userID}", wrapJSON(api.UserGet)).
		Methods("GET")
	router.HandleFunc("/users", wrapJSON(api.UserCreate)).
		Methods("POST")
	router.HandleFunc("/users/{userID}", wrapJSON(api.UserUpdate)).
		Methods("PUT")
	router.HandleFunc("/users/{userID}", wrapJSON(api.UserDelete)).
		Methods("DELETE")

	router.HandleFunc("/twitterAccounts", wrapJSON(api.TwitterAccountsAll)).
		Methods("GET")
	router.HandleFunc("/twitterAccounts/{twitterAccountID}", wrapJSON(api.TwitterAccountGet)).
		Methods("GET")
	router.HandleFunc("/twitterAccounts/{twitterAccountID}/tweets", wrapJSON(api.TwitterAccountGetWithTweets)).
		Methods("GET")

	server := http.Server{
		Addr:    *addr,
		Handler: router,
	}

	log.Printf("GoBot Data Server running on %s...\n", *addr)
	server.ListenAndServe()
}

func wrapJSON(apiFunc func(http.ResponseWriter, *http.Request) interface{}) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		response := apiFunc(res, req)

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
