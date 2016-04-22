package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

// TwitterAccountsAll = GET: /twitterAccounts
func TwitterAccountsAll(res http.ResponseWriter, req *http.Request) {
	qs := req.URL.Query()

	recordsPerPage, err := strconv.Atoi(qs.Get("recordsPerPage"))
	if err != nil {
		recordsPerPage = 20
	}
	if recordsPerPage < 1 {
		recordsPerPage = 1
	} else if recordsPerPage > 100 {
		recordsPerPage = 100
	}

	page, err := strconv.Atoi(qs.Get("page"))
	if err != nil {
		page = 1
	}

	if page < 1 {
		page = 1
	} else if page > 100 {
		page = 100
	}

	paging := db.PagingInfo{
		OrderBy: db.UsersOrderByDateCreated,
		Asc:     false,
		Limit:   page * recordsPerPage,
		Offset:  (page - 1) * recordsPerPage,
	}

	twitterAccounts, err := db.TwitterAccountsAll(paging)
	if err != nil {
		panic(err)
	}

	accounts := make([]TwitterAccount, len(twitterAccounts))

	for _, twitterAccount := range twitterAccounts {
		account := TwitterAccount{
			ID:                twitterAccount.ID,
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
				ID:       tweet.ID,
				Text:     tweet.Tweet,
				PostOn:   tweet.PostOn,
				IsPosted: tweet.IsPosted,
			})
		}

		accounts = append(accounts, account)
	}

	model := struct {
		Message         string           `json:"message"`
		TwitterAccounts []TwitterAccount `json:"twitterAccounts"`
	}{}

	model.Message = ok
	model.TwitterAccounts = accounts

	data, err := json.MarshalIndent(model, "", "    ")
	if err != nil {
		panic(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}

// TwitterAccountGet = GET: /twitterAccount/{twitterAccountID}
func TwitterAccountGet(res http.ResponseWriter, req *http.Request) {
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
		ID:                account.ID,
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
			ID:       tweet.ID,
			Text:     tweet.Tweet,
			PostOn:   tweet.PostOn,
			IsPosted: tweet.IsPosted,
		})
	}
}
