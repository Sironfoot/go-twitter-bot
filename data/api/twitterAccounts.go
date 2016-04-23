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

// TwitterAccountsAll = GET: /twitterAccounts
func TwitterAccountsAll(res http.ResponseWriter, req *http.Request) {
	var response interface{}

	defer func() {
		res.Header().Set("Content-Type", "application/json")

		data, jsonErr := json.MarshalIndent(response, "", "    ")
		if jsonErr != nil {
			panic(jsonErr)
		}
		res.Write(data)
	}()

	paging, errResponse := extractAndValidatePagingInfo(req)
	if errResponse != nil {
		response = errResponse
		return
	}

	twitterAccounts, totalRecords, err := db.TwitterAccountsAll(paging)
	if err != nil {
		panic(err)
	}

	var accounts []TwitterAccount

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

		accounts = append(accounts, account)
	}

	model := struct {
		pagedResponse
		TwitterAccounts []TwitterAccount `json:"twitterAccounts"`
	}{}

	model.Message = ok
	model.Page = paging.Page
	model.RecordsPerPage = paging.RecordsPerPage
	model.TotalRecords = totalRecords

	if len(accounts) == 0 {
		model.TwitterAccounts = make([]TwitterAccount, 0)
	} else {
		model.TwitterAccounts = accounts
	}

	response = model
}

// TwitterAccountGet = GET: /twitterAccount/{twitterAccountID}
func TwitterAccountGet(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	twitterAccountID := vars["twitterAccountID"]

	model := struct {
		messageResponse
		TwitterAccount TwitterAccount `json:"twitterAccount"`
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
