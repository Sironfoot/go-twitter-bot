package db

import (
	"time"

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
	AccessTokenSecret string    `db:"accessTokenSecret"`
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
var TwitterAccountFromID = func(id string) (TwitterAccount, error) {
	var account TwitterAccount

	if !isUUID.MatchString(id) {
		return account, ErrEntityNotFound
	}

	err := sqlboiler.EntityGetByID(&account, id, db)
	if err == sqlboiler.ErrEntityNotFound {
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

// TwitterAccountsAll returns all TwitterAccount records from the database
var TwitterAccountsAll = func(paging PagingInfo) ([]TwitterAccount, error) {
	var accounts []TwitterAccount

	cmd := `SELECT id, ` + sqlboiler.GetColumnListString(&TwitterAccount{}) + `
			FROM twitter_accounts
			ORDER BY $1
			LIMIT $2 OFFSET $3`

	rows, err := dbx.Queryx(cmd, paging.BuildOrderBy(), paging.Limit, paging.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		account := TwitterAccount{}
		err := rows.StructScan(&account)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

// TwitterAccountGetTweets loads Tweets child entites for TwitterAccount
var TwitterAccountGetTweets = func(account *TwitterAccount) ([]Tweet, error) {
	var tweets []Tweet

	if !account.IsTransient() {
		cmd := `SELECT id, tweet, post_on, is_posted, date_created
				FROM tweets
				WHERE twitter_account_id = $1`

		rows, err := db.Query(cmd, account.ID)
		if err != nil {
			return tweets, err
		}

		for rows.Next() {
			tweet := Tweet{AccountID: account.ID}
			err = rows.Scan(
				&tweet.ID,
				&tweet.Tweet,
				&tweet.PostOn,
				&tweet.IsPosted,
				&tweet.DateCreated)

			if err != nil {
				return tweets, err
			}

			tweets = append(tweets, tweet)
		}

		rows.Close()
	}

	return tweets, nil
}

// GetTweets loads Tweets child entites for TwitterAccount
func (account *TwitterAccount) GetTweets() ([]Tweet, error) {
	return TwitterAccountGetTweets(account)
}
