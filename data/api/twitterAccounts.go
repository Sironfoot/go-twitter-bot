package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sironfoot/go-twitter-bot/data/db"
)

// TwitterAccount model returned by REST API
type TwitterAccount struct {
	ID                string    `json:"id"`
	UserID            string    `json:"userId"`
	Username          string    `json:"username"`
	DateCreated       time.Time `json:"dateCreated"`
	ConsumerKey       string    `json:"consumerKey"`
	ConsumerSecret    string    `json:"consumerSecret"`
	AccessToken       string    `json:"accessToken"`
	AccessTokenSecret string    `json:"accessTokenSecret"`

	Tweets []Tweet `json:"tweets"`
}

// Tweet model returned by REST API
type Tweet struct {
	ID       string    `json:"id"`
	Text     string    `json:"text"`
	PostOn   time.Time `json:"postOn"`
	IsPosted bool      `json:"isPosted"`
}

// GetTwitterAccounts = GET: /twitterAccounts
func GetTwitterAccounts(res http.ResponseWriter, req *http.Request) {
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

	data, err := json.MarshalIndent(accounts, "", "    ")
	if err != nil {
		panic(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}

// GetTwitterAccount = GET: /twitterAccount/{twitterAccountID}
func GetTwitterAccount(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	twitterAccountID := vars["twitterAccountID"]

	model := struct {
		Message        string         `json:"message"`
		TwitterAccount TwitterAccount `json:"TwitterAccount"`
	}{}

	res.Header().Set("Content-Type", "application/json")

	defer func() {
		data, err := json.MarshalIndent(model, "", "    ")
		if err != nil {
			panic(err)
		}

		res.Write(data)
	}()

	account, err := db.TwitterAccountFromID(twitterAccountID)
	if err == db.ErrEntityNotFound {
		model.Message = fmt.Sprintf("TwitterAccount not found on ID: %s", twitterAccountID)
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		panic(err)
	}

	model.Message = "OK"
	model.TwitterAccount = TwitterAccount{
		ID:                account.ID(),
		UserID:            account.UserID,
		Username:          account.Username,
		DateCreated:       account.DateCreated,
		ConsumerKey:       account.ConsumerKey,
		ConsumerSecret:    account.ConsumerSecret,
		AccessToken:       account.AccessToken,
		AccessTokenSecret: account.AccessTokenSecret,
	}

	tweets, err := account.GetTweets()
	if err != nil {
		panic(err)
	}

	for _, tweet := range tweets {
		model.TwitterAccount.Tweets = append(model.TwitterAccount.Tweets, Tweet{
			ID:       tweet.ID(),
			Text:     tweet.Tweet,
			PostOn:   tweet.PostOn,
			IsPosted: tweet.IsPosted,
		})
	}
}
