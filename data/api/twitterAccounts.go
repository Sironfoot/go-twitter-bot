package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"goji.io/pat"

	"golang.org/x/net/context"

	"github.com/sironfoot/go-twitter-bot/data/db"
	"github.com/sironfoot/go-twitter-bot/data/models"
)

type twitterAccountBase struct {
	ID                string    `json:"id"`
	UserID            string    `json:"userId"`
	Username          string    `json:"username"`
	DateCreated       time.Time `json:"dateCreated"`
	ConsumerKey       string    `json:"consumerKey"`
	ConsumerSecret    string    `json:"consumerSecret"`
	AccessToken       string    `json:"accessToken"`
	AccessTokenSecret string    `json:"accessTokenSecret"`
}

type twitterAccount struct {
	twitterAccountBase
	Tweets int `json:"tweets"`
}

type twitterAccountWithTweets struct {
	twitterAccountBase
	Tweets childTweets `json:"tweets"`
}

type childTweets struct {
	Page           int     `json:"page"`
	RecordsPerPage int     `json:"recordPerPage"`
	TotalRecords   int     `json:"totalRecords"`
	Records        []tweet `json:"records"`
}

type tweet struct {
	ID       string    `json:"id"`
	Text     string    `json:"text"`
	PostOn   time.Time `json:"postOn"`
	IsPosted bool      `json:"isPosted"`
}

// TwitterAccountsAll = GET: /twitterAccounts
func TwitterAccountsAll(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appContext := ctx.Value("appContext").(*AppContext)

	defaults := getPagingDefaults(db.TwitterAccountsOrderByDateCreated, false, db.TwitterAccountsSortableColumns)
	paging, err := ExtractAndValidatePagingInfo(req, defaults)
	if err != nil {
		appContext.Response = MessageResponse{err.Error()}
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

	filterUserID := qs.Get("userID")
	// filter TwitterAccounts to user's own if not an admin
	if !appContext.AuthUser.IsAdmin {
		filterUserID = appContext.AuthUser.ID
	}
	query.UserID = filterUserID // TODO: implement UserID support in data access code

	twitterAccounts, totalRecords, err := db.TwitterAccountsAll(query)
	if err != nil {
		panic(err)
	}

	var accounts []twitterAccount

	for _, accountDB := range twitterAccounts {
		account := twitterAccount{
			twitterAccountBase: twitterAccountBase{
				ID:                accountDB.ID,
				UserID:            accountDB.UserID,
				Username:          accountDB.Username,
				DateCreated:       accountDB.DateCreated,
				ConsumerKey:       accountDB.ConsumerKey,
				ConsumerSecret:    accountDB.ConsumerSecret,
				AccessToken:       accountDB.AccessToken,
				AccessTokenSecret: accountDB.AccessTokenSecret,
			},
			Tweets: accountDB.NumTweets,
		}

		accounts = append(accounts, account)
	}

	model := struct {
		pagedResponse
		TwitterAccounts []twitterAccount `json:"twitterAccounts"`
	}{}

	model.Message = ok
	model.Page = paging.Page
	model.RecordsPerPage = paging.RecordsPerPage
	model.TotalRecords = totalRecords

	if len(accounts) == 0 {
		model.TwitterAccounts = make([]twitterAccount, 0)
	} else {
		model.TwitterAccounts = accounts
	}

	appContext.Response = model
}

// TwitterAccountGet = GET: /twitterAccount/:twitterAccountID
func TwitterAccountGet(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appContext := ctx.Value("appContext").(*AppContext)
	twitterAccountID := pat.Param(ctx, "twitterAccountID")

	account, err := db.TwitterAccountFromID(twitterAccountID)
	if err == db.ErrEntityNotFound {
		res.WriteHeader(http.StatusNotFound)
		appContext.Response = MessageResponse{
			Message: fmt.Sprintf("TwitterAccount not found on ID: %s", twitterAccountID),
		}
		return
	} else if err != nil {
		panic(err)
	}

	// non-admins can only view their own TwitterAccounts
	if !appContext.AuthUser.IsAdmin && appContext.AuthUser.ID != account.UserID {
		appContext.Response = MessageResponse{
			Message: "This resource is only available to users with administrator rights.",
		}
		res.WriteHeader(http.StatusForbidden)
		return
	}

	// non-admins can only view their own TwitterAccounts
	if !appContext.AuthUser.IsAdmin && account.UserID != appContext.AuthUser.ID {
		appContext.Response = MessageResponse{
			Message: "This resource is only available to users with administrator rights.",
		}
		res.WriteHeader(http.StatusForbidden)
		return
	}

	model := struct {
		MessageResponse
		TwitterAccount twitterAccount `json:"twitterAccount"`
	}{}

	model.Message = "OK"
	model.TwitterAccount = twitterAccount{
		twitterAccountBase: twitterAccountBase{
			ID:                account.ID,
			UserID:            account.UserID,
			Username:          account.Username,
			DateCreated:       account.DateCreated,
			ConsumerKey:       account.ConsumerKey,
			ConsumerSecret:    account.ConsumerSecret,
			AccessToken:       account.AccessToken,
			AccessTokenSecret: account.AccessTokenSecret,
		},
		Tweets: account.NumTweets,
	}

	appContext.Response = model
}

// TwitterAccountGetWithTweets = GET: /twitterAccounts/:twitterAccountID/tweets
func TwitterAccountGetWithTweets(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appContext := ctx.Value("appContext").(*AppContext)
	twitterAccountID := pat.Param(ctx, "twitterAccountID")

	account, err := db.TwitterAccountFromID(twitterAccountID)
	if err == db.ErrEntityNotFound {
		res.WriteHeader(http.StatusNotFound)
		appContext.Response = MessageResponse{
			Message: fmt.Sprintf("TwitterAccount not found on ID: %s", twitterAccountID),
		}
		return
	} else if err != nil {
		panic(err)
	}

	// non-admins can only view their own TwitterAccounts
	if !appContext.AuthUser.IsAdmin && appContext.AuthUser.ID != account.UserID {
		appContext.Response = MessageResponse{
			Message: "This resource is only available to users with administrator rights.",
		}
		res.WriteHeader(http.StatusForbidden)
		return
	}

	model := struct {
		MessageResponse
		TwitterAccount twitterAccountWithTweets `json:"twitterAccount"`
	}{}

	model.Message = "OK"
	model.TwitterAccount = twitterAccountWithTweets{
		twitterAccountBase: twitterAccountBase{
			ID:                account.ID,
			UserID:            account.UserID,
			Username:          account.Username,
			DateCreated:       account.DateCreated,
			ConsumerKey:       account.ConsumerKey,
			ConsumerSecret:    account.ConsumerSecret,
			AccessToken:       account.AccessToken,
			AccessTokenSecret: account.AccessTokenSecret,
		},
		Tweets: childTweets{
			Page:           1,
			TotalRecords:   0,
			RecordsPerPage: 20,
		},
	}

	defaults := getPagingDefaults(db.TweetsOrderByDateCreated, false, db.TweetsSortableColumns)
	paging, err := ExtractAndValidatePagingInfo(req, defaults)
	if err != nil {
		appContext.Response = MessageResponse{err.Error()}
		return
	}

	query := db.TweetsQuery{
		PagingInfo: paging,
	}

	qs := req.URL.Query()
	dateTime, err := time.Parse("2006-01-02 15:04:05", qs.Get("tweetsToBePostedSince"))
	if err == nil {
		query.ToBePostedSince = dateTime
	}

	tweets, totalTweets, err := account.GetTweets(query)
	if err != nil {
		panic(err)
	}

	model.TwitterAccount.Tweets.TotalRecords = totalTweets

	if len(tweets) > 0 {
		for _, tweetDB := range tweets {
			model.TwitterAccount.Tweets.Records = append(model.TwitterAccount.Tweets.Records, tweet{
				ID:       tweetDB.ID,
				Text:     tweetDB.Tweet,
				PostOn:   tweetDB.PostOn,
				IsPosted: tweetDB.IsPosted,
			})
		}
	} else {
		model.TwitterAccount.Tweets.Records = make([]tweet, 0)
	}

	appContext.Response = model
}

// TwitterAccountTweetCreate = POST: /twitterAccounts/:twitterAccountID/tweets
func TwitterAccountTweetCreate(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appContext := ctx.Value("appContext").(*AppContext)
	twitterAccountID := pat.Param(ctx, "twitterAccountID")

	account, err := db.TwitterAccountFromID(twitterAccountID)
	if err == db.ErrEntityNotFound {
		res.WriteHeader(http.StatusNotFound)
		appContext.Response = MessageResponse{
			Message: fmt.Sprintf("TwitterAccount not found on ID: %s", twitterAccountID),
		}
		return
	} else if err != nil {
		panic(err)
	}

	// non-admins can only create Tweets for their own TwitterAccounts
	if !appContext.AuthUser.IsAdmin && appContext.AuthUser.ID != account.UserID {
		appContext.Response = MessageResponse{
			Message: "This resource is only available to users with administrator rights.",
		}
		res.WriteHeader(http.StatusForbidden)
		return
	}

	var newTweet models.Tweet

	err = json.NewDecoder(req.Body).Decode(&newTweet)
	if err != nil {
		panic(err)
	}
	req.Body.Close()

	newTweet.Sanitise()
	validationErrors, err := newTweet.ValidateCreate()
	if err != nil {
		panic(err)
	}

	model := createResponse{}

	if len(validationErrors) > 0 {
		model.Message = "Tweet model is invalid."
		model.Errors = validationErrors
		appContext.Response = model

		res.WriteHeader(http.StatusBadRequest)
		return
	}

	tweet := &db.Tweet{
		AccountID:   account.ID,
		Tweet:       newTweet.Text,
		PostOn:      newTweet.PostOn,
		IsPosted:    newTweet.IsPosted,
		DateCreated: time.Now().UTC(),
	}

	err = tweet.Save()
	if err != nil {
		panic(err)
	}

	model.Message = ok
	model.ID = &tweet.ID
	res.WriteHeader(http.StatusCreated)

	appContext.Response = model
}

// TwitterAccountTweetUpdate = PUT: /twitterAccounts/:twitterAccountID/tweets/:tweetID
func TwitterAccountTweetUpdate(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appContext := ctx.Value("appContext").(*AppContext)
	twitterAccountID := pat.Param(ctx, "twitterAccountID")

	account, err := db.TwitterAccountFromID(twitterAccountID)
	if err == db.ErrEntityNotFound {
		res.WriteHeader(http.StatusNotFound)
		appContext.Response = MessageResponse{
			Message: fmt.Sprintf("TwitterAccount not found on ID: %s", twitterAccountID),
		}
		return
	} else if err != nil {
		panic(err)
	}

	// non-admins can only edit their own Tweets
	if !appContext.AuthUser.IsAdmin && appContext.AuthUser.ID != account.UserID {
		appContext.Response = MessageResponse{
			Message: "This resource is only available to users with administrator rights.",
		}
		res.WriteHeader(http.StatusForbidden)
		return
	}

	tweetID := pat.Param(ctx, "tweetID")
	tweet, err := account.GetTweetFromID(tweetID)
	if err == db.ErrEntityNotFound {
		res.WriteHeader(http.StatusNotFound)
		appContext.Response = MessageResponse{
			Message: fmt.Sprintf("Tweet not found on ID: %s", tweetID),
		}
		return
	} else if err != nil {
		panic(err)
	}

	var updateTweet models.Tweet

	err = json.NewDecoder(req.Body).Decode(&updateTweet)
	if err != nil {
		panic(err)
	}
	req.Body.Close()

	updateTweet.Sanitise()
	validationErrors, err := updateTweet.ValidateUpdate(tweetID)

	if len(validationErrors) > 0 {
		res.WriteHeader(http.StatusBadRequest)
		appContext.Response = updateResponse{
			Message: "Tweet model is invalid.",
			Errors:  validationErrors,
		}
		return
	}

	tweet.Tweet = updateTweet.Text
	tweet.PostOn = updateTweet.PostOn
	tweet.IsPosted = updateTweet.IsPosted

	err = tweet.Save()
	if err != nil {
		panic(err)
	}

	appContext.Response = MessageResponse{
		Message: ok,
	}
}

// TwitterAccountTweetDelete = DELETE: /twitterAccounts/:twitterAccountID/tweets/:tweetID
func TwitterAccountTweetDelete(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appContext := ctx.Value("appContext").(*AppContext)
	twitterAccountID := pat.Param(ctx, "twitterAccountID")

	account, err := db.TwitterAccountFromID(twitterAccountID)
	if err == db.ErrEntityNotFound {
		res.WriteHeader(http.StatusNotFound)
		appContext.Response = MessageResponse{
			Message: fmt.Sprintf("TwitterAccount not found on ID: %s", twitterAccountID),
		}
		return
	} else if err != nil {
		panic(err)
	}

	// non-admins can only delete their own Tweets
	if !appContext.AuthUser.IsAdmin && appContext.AuthUser.ID != account.UserID {
		appContext.Response = MessageResponse{
			Message: "This resource is only available to users with administrator rights.",
		}
		res.WriteHeader(http.StatusForbidden)
		return
	}

	tweetID := pat.Param(ctx, "tweetID")
	tweet, err := account.GetTweetFromID(tweetID)
	if err == db.ErrEntityNotFound {
		res.WriteHeader(http.StatusNotFound)
		appContext.Response = MessageResponse{
			Message: fmt.Sprintf("Tweet not found on ID: %s", tweetID),
		}
		return
	} else if err != nil {
		panic(err)
	}

	err = tweet.Delete()
	if err != nil {
		panic(err)
	}

	appContext.Response = MessageResponse{
		Message: ok,
	}
}
