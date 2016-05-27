package db

import (
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/sironfoot/go-twitter-bot/lib/sqlboiler"
)

// TwitterAccount maps to twitter_accounts table
type TwitterAccount struct {
	ID                string    `db:"id"`
	UserID            string    `db:"user_id"`
	Username          string    `db:"username"`
	DateCreated       time.Time `db:"date_created"`
	ConsumerKey       string    `db:"consumer_key"`
	ConsumerSecret    string    `db:"consumer_secret"`
	AccessToken       string    `db:"access_token"`
	AccessTokenSecret string    `db:"access_token_secret"`
}

// IsTransient determines if TwitterAccount record has been saved to the database,
// true means TwitterAccount struct has NOT been saved, false means it has.
func (account *TwitterAccount) IsTransient() bool {
	return len(account.ID) == 0
}

// MetaData returns meta data information about the TwitterAccount entity
func (account *TwitterAccount) MetaData() sqlboiler.EntityMetaData {
	return sqlboiler.EntityMetaData{
		TableName:      "twitter_accounts",
		PrimaryKeyName: "id",
	}
}

// TwitterAccountSave saves the TwitterAccount struct to the database.
var TwitterAccountSave = func(account *TwitterAccount) error {
	return sqlboiler.EntitySave(account, db)
}

// Save saves the TwitterAccount struct to the database.
func (account *TwitterAccount) Save() error {
	return TwitterAccountSave(account)
}

// TwitterAccountDelete deletes the TwitterAccount from the database
var TwitterAccountDelete = func(account *TwitterAccount) error {
	return sqlboiler.EntityDelete(account, db)
}

// Delete deletes the TwitterAccount from the database
func (account *TwitterAccount) Delete() error {
	return TwitterAccountDelete(account)
}

// TwitterAccountFromID returns a TwitterAccount record with given ID
var TwitterAccountFromID = func(id string) (TwitterAccountList, error) {
	var account TwitterAccountList

	if !isUUID.MatchString(id) {
		return account, ErrEntityNotFound
	}

	cmd := `SELECT ` + sqlboiler.GetFullColumnListString(&TwitterAccount{}, "ta") + `, COUNT(t.id) AS num_tweets
			FROM twitter_accounts ta
				LEFT OUTER JOIN tweets t ON ta.id = t.twitter_account_id
			WHERE ta.id = $1
			GROUP BY ta.id`

	err := dbx.QueryRowx(cmd, id).StructScan(&account)
	if err == sql.ErrNoRows {
		return account, ErrEntityNotFound
	}
	return account, err
}

const (
	// TwitterAccountsOrderByUsername is for ordering TwitterAccounts by Username
	TwitterAccountsOrderByUsername = "username"
	// TwitterAccountsOrderByDateCreated is for ordering TwitterAccounts by DateCreated
	TwitterAccountsOrderByDateCreated = "date_created"
)

// TwitterAccountsSortableColumns is a list of allowed sortable columns
var TwitterAccountsSortableColumns = []string{
	TwitterAccountsOrderByUsername,
	TwitterAccountsOrderByDateCreated,
}

// TwitterAccountList is a TwitterAccount struct that includes a NumTweets field
type TwitterAccountList struct {
	TwitterAccount
	NumTweets int `db:"num_tweets"`
}

// TwitterAccountQuery is a search query for searching TwitterAccounts, used by db.TwitterAccountsAll
type TwitterAccountQuery struct {
	PagingInfo
	ContainsUsername         string
	UserID                   string
	HasTweetsToBePostedSince time.Time
}

// TwitterAccountsAll returns all TwitterAccount records from the database
var TwitterAccountsAll = func(query TwitterAccountQuery) ([]TwitterAccountList, int, error) {
	var accounts []TwitterAccountList
	recordCount := 0

	cmd := sq.
		Select(sqlboiler.GetFullColumnList(&TwitterAccount{}, "ta")...).Column("COUNT(t.id) AS num_tweets").
		From("twitter_accounts ta").
		LeftJoin("tweets t ON ta.id = t.twitter_account_id")

	if query.ContainsUsername != "" {
		cmd = cmd.Where("ta.username LIKE ?", query.ContainsUsername)
	}

	if !query.HasTweetsToBePostedSince.IsZero() {
		cmd = cmd.Where("t.is_posted = ? AND t.post_on > ?", false, query.HasTweetsToBePostedSince)
	}

	if query.UserID != "" {
		cmd = cmd.Where("ta.user_id = ?", query.UserID)
	}

	cmd = cmd.
		GroupBy("ta.id").
		OrderBy(query.BuildOrderBy()).
		Limit(uint64(query.Limit())).
		Offset(uint64(query.Offset()))

	sqlString, args, err := cmd.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, recordCount, err
	}

	fmt.Println(sqlString)
	fmt.Println(args)

	rows, err := dbx.Queryx(sqlString, args...)
	if err != nil {
		return nil, recordCount, err
	}

	defer rows.Close()

	for rows.Next() {
		account := TwitterAccountList{}
		err = rows.StructScan(&account)

		if err != nil {
			return nil, recordCount, err
		}

		accounts = append(accounts, account)
	}

	countCmd := sq.
		Select("COUNT(DISTINCT ta.id)").
		From("twitter_accounts ta")

	if query.ContainsUsername != "" || !query.HasTweetsToBePostedSince.IsZero() {
		countCmd = countCmd.
			LeftJoin("tweets t ON ta.id = t.twitter_account_id")

		if query.ContainsUsername != "" {
			countCmd = countCmd.Where("ta.username LIKE ?", query.ContainsUsername)
		}

		if !query.HasTweetsToBePostedSince.IsZero() {
			countCmd = countCmd.Where("t.is_posted = ? AND t.post_on > ?", false, query.HasTweetsToBePostedSince)
		}

		if query.UserID != "" {
			countCmd = countCmd.Where("ta.user_id = ?", query.UserID)
		}
	}

	countSQL, countArgs, err := countCmd.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, recordCount, err
	}

	err = dbx.Get(&recordCount, countSQL, countArgs...)

	return accounts, recordCount, err
}

const (
	// TweetsOrderByDateCreated is for ordering Tweets by DateCreated
	TweetsOrderByDateCreated = "date_created"
)

// TweetsSortableColumns is a list of allowed sortable columns
var TweetsSortableColumns = []string{
	TweetsOrderByDateCreated,
}

// TweetsQuery is a search query for searching Tweets, used by db.TwitterAccountGetTweets
type TweetsQuery struct {
	PagingInfo
	ToBePostedSince time.Time
}

// TwitterAccountGetTweets loads Tweets child entites for TwitterAccount
var TwitterAccountGetTweets = func(account *TwitterAccount, query TweetsQuery) ([]Tweet, int, error) {
	var tweets []Tweet
	totalRecords := 0

	if account.IsTransient() {
		return tweets, totalRecords, nil
	}

	var queryParams []interface{}
	queryParams = append(queryParams, account.ID)
	queryParams = append(queryParams, query.BuildOrderBy())
	queryParams = append(queryParams, query.Limit())
	queryParams = append(queryParams, query.Offset())

	cmd := `SELECT id, ` + sqlboiler.GetColumnListString(&Tweet{}, "") + `
			FROM tweets
			WHERE twitter_account_id = $1 `

	if !query.ToBePostedSince.IsZero() {
		cmd += `AND is_posted = false AND post_on > $5 `
		queryParams = append(queryParams, query.ToBePostedSince)
	}

	cmd += `ORDER BY $2
			LIMIT $3 OFFSET $4`

	rows, err := dbx.Queryx(cmd, queryParams...)
	if err != nil {
		return tweets, totalRecords, err
	}

	for rows.Next() {
		tweet := Tweet{}
		err = rows.StructScan(&tweet)

		if err != nil {
			return tweets, totalRecords, err
		}

		tweets = append(tweets, tweet)
	}

	rows.Close()

	// count tweets
	var countParams []interface{}
	countCmd := `SELECT COUNT(*) FROM tweets WHERE twitter_account_id = $1`

	countParams = append(countParams, account.ID)

	if !query.ToBePostedSince.IsZero() {
		countCmd += ` AND is_posted = false AND post_on > $2`
		countParams = append(countParams, query.ToBePostedSince)
	}

	err = dbx.Get(&totalRecords, countCmd, countParams...)

	return tweets, totalRecords, err
}

// GetTweets loads Tweets child entites for TwitterAccount
func (account *TwitterAccount) GetTweets(query TweetsQuery) ([]Tweet, int, error) {
	return TwitterAccountGetTweets(account, query)
}

// TwitterAccountGetTweetFromID gets a TwitterAccount's Tweet by its ID
var TwitterAccountGetTweetFromID = func(account *TwitterAccount, tweetID string) (Tweet, error) {
	var tweet Tweet

	cmd := `SELECT id, ` + sqlboiler.GetColumnListString(&tweet, "") + `
			FROM tweets
			WHERE twitter_account_id = $1 AND id = $2`

	err := dbx.Get(&tweet, cmd, account.ID, tweetID)
	if err == sql.ErrNoRows {
		return tweet, ErrEntityNotFound
	}
	return tweet, err
}

// GetTweetFromID gets this TwitterAccount's Tweet by ID
func (account *TwitterAccount) GetTweetFromID(id string) (Tweet, error) {
	return TwitterAccountGetTweetFromID(account, id)
}
