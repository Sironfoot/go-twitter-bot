package api

import (
	"encoding/json"
	"net/http"
	"time"

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
