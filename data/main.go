package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sironfoot/go-twitter-bot/data/db"
)

var addr = flag.String("addr", "localhost:7001", "Address to run server on")

// TwitterAccount model returned by REST API
type TwitterAccount struct {
	ID                string
	UserID            string
	Username          string
	DateCreated       time.Time
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string

	Tweets []Tweet
}

// Tweet model returned by REST API
type Tweet struct {
	ID       string
	Text     string
	PostOn   time.Time
	IsPosted bool
}

func main() {
	flag.Parse()
	db.InitDB("user=postgres dbname=go_twitter_bot sslmode=disable")

	router := mux.NewRouter()

	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Hello from GoBot Data server\n")
	})

	router.HandleFunc("/users", func(res http.ResponseWriter, req *http.Request) {
		users, err := db.UsersAll()
		if err != nil {
			panic(err)
		}

		data, err := json.MarshalIndent(users, "", "\t")
		if err != nil {
			panic(err)
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
	})

	router.HandleFunc("/twitterAccounts", func(res http.ResponseWriter, req *http.Request) {
		twitterAccounts, err := db.TwitterAccountsAll()
		if err != nil {
			panic(err)
		}

		var accounts []TwitterAccount

		for _, twitterAccount := range twitterAccounts {
			account := TwitterAccount{
				ID:                twitterAccount.ID(),
				UserID:            twitterAccount.UserID,
				Username:          twitterAccount.Username,
				DateCreated:       twitterAccount.DateCreated,
				ConsumerKey:       twitterAccount.ConsumerKey,
				ConsumerSecret:    twitterAccount.ConsumerSecret,
				AccessToken:       twitterAccount.AccessToken,
				AccessTokenSecret: twitterAccount.AccessTokenSecret,
			}

			tweets, tweetErr := twitterAccount.GetTweets()
			if tweetErr != nil {
				panic(tweetErr)
			}

			for _, tweet := range tweets {
				account.Tweets = append(account.Tweets, Tweet{
					ID:       tweet.ID(),
					Text:     tweet.Tweet,
					PostOn:   tweet.PostOn,
					IsPosted: tweet.IsPosted,
				})
			}

			accounts = append(accounts, account)
		}

		data, err := json.MarshalIndent(accounts, "", "\t")
		if err != nil {
			panic(err)
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
	})

	server := http.Server{
		Addr:    *addr,
		Handler: router,
	}

	log.Printf("GoBot Data Server running on %s...\n", *addr)
	server.ListenAndServe()
}
