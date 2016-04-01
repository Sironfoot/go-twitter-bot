package main

import (
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

	router.HandleFunc("/users", api.GetUsers).
		Methods("GET")
	router.HandleFunc("/users", api.CreateUser).
		Methods("POST")
	router.HandleFunc("/users/{userID}", api.GetUser).
		Methods("GET")
	router.HandleFunc("/users/{userID}", api.UpdateUser).
		Methods("PUT")

	router.HandleFunc("/twitterAccounts", api.GetTwitterAccounts).
		Methods("GET")
	router.HandleFunc("/twitterAccounts/{twitterAccountID}", api.GetTwitterAccount).
		Methods("GET")

	server := http.Server{
		Addr:    *addr,
		Handler: router,
	}

	log.Printf("GoBot Data Server running on %s...\n", *addr)
	server.ListenAndServe()
}
