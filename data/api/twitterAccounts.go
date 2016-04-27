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
	NumTweets         int       `json:"numTweets"`
}

// TwitterAccountWithTweets model returned by GetByID API
type TwitterAccountWithTweets struct {
	TwitterAccount
	Tweets ChildTweets `json:"tweets"`
}

// ChildTweets ...
type ChildTweets struct {
	Page           int     `json:"page"`
	RecordsPerPage int     `json:"recordPerPage"`
	TotalRecords   int     `json:"totalRecords"`
	Records        []Tweet `json:"records"`
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

	query := db.TwitterAccountQuery{
		PagingInfo: paging,
	}

	qs := req.URL.Query()
	query.ContainsUsername = qs.Get("username")
	dateTime, err := time.Parse("2006-01-02 15:04:05", qs.Get("hasTweetsToBePostedSince"))
	if err == nil {
		query.HasTweetsToBePostedSince = dateTime
	}

	twitterAccounts, totalRecords, err := db.TwitterAccountsAll(query)
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
			NumTweets:         twitterAccount.NumTweets,
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
		TwitterAccount TwitterAccountWithTweets `json:"twitterAccount"`
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
	model.TwitterAccount = TwitterAccountWithTweets{
		TwitterAccount: TwitterAccount{
			ID:                account.ID,
			UserID:            account.UserID,
			Username:          account.Username,
			DateCreated:       account.DateCreated,
			ConsumerKey:       account.ConsumerKey,
			ConsumerSecret:    account.ConsumerSecret,
			AccessToken:       account.AccessToken,
			AccessTokenSecret: account.AccessTokenSecret,
			NumTweets:         account.NumTweets,
		},
	}

	tweets, err := account.GetTweets()
	if err != nil {
		panic(err)
	}

	model.TwitterAccount.Tweets = ChildTweets{
		Page:           1,
		TotalRecords:   2,
		RecordsPerPage: 20,
	}

	for _, tweet := range tweets {
		model.TwitterAccount.Tweets.Records = append(model.TwitterAccount.Tweets.Records, Tweet{
			ID:       tweet.ID,
			Text:     tweet.Tweet,
			PostOn:   tweet.PostOn,
			IsPosted: tweet.IsPosted,
		})
	}
}
