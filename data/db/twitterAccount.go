package db

import (
	"database/sql"
	"time"
)

// TwitterAccount maps to twitter_accounts table
type TwitterAccount struct {
	ID                string
	UserID            string
	Username          string
	DateCreated       time.Time
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// IsTransient determines if TwitterAccount record has been saved to the database,
// true means TwitterAccount struct has NOT been saved, false means it has.
func (account *TwitterAccount) IsTransient() bool {
	return len(account.ID) == 0
}

// TwitterAccountSave saves the TwitterAccount struct to the database.
var TwitterAccountSave = func(account *TwitterAccount) error {
	if account.IsTransient() {
		cmd := `INSERT INTO twitter_accounts
				(
					user_id,
					username,
					date_created,
					consumer_key,
					consumer_secret,
					access_token,
					access_token_secret
				)
				VALUES($1, $2, $3, $4, $5, $6, $7)
				RETURNING id`

		statement, err := db.Prepare(cmd)
		if err != nil {
			return err
		}
		defer statement.Close()

		err = statement.
			QueryRow(
				account.UserID,
				account.Username,
				account.DateCreated,
				account.ConsumerKey,
				account.ConsumerSecret,
				account.AccessToken,
				account.AccessTokenSecret).
			Scan(&account.ID)
		if err != nil {
			return err
		}
	} else {
		cmd := `UPDATE twitter_accounts
				SET user_id = $2,
					username = $3,
					date_created = $4,
					consumer_key = $5,
					consumer_secret = $6,
					access_token = $7,
					access_token_secret = $8
				WHERE id = $1`

		_, err := db.Exec(cmd,
			account.ID,
			account.UserID,
			account.Username,
			account.DateCreated,
			account.ConsumerKey,
			account.ConsumerSecret,
			account.AccessToken,
			account.AccessTokenSecret)
		if err != nil {
			return err
		}
	}

	return nil
}

// Save saves the TwitterAccount struct to the database.
func (account *TwitterAccount) Save() error {
	return TwitterAccountSave(account)
}

// TwitterAccountDelete deletes the TwitterAccount from the database
var TwitterAccountDelete = func(account *TwitterAccount) error {
	cmd := `DELETE FROM twitter_accounts
			WHERE id = $1`

	_, err := db.Exec(cmd, account.ID)
	return err
}

// Delete deletes the TwitterAccount from the database
func (account *TwitterAccount) Delete() error {
	return TwitterAccountDelete(account)
}

// TwitterAccountFromID returns a TwitterAccount record with given ID
var TwitterAccountFromID = func(id string) (TwitterAccount, error) {
	var account TwitterAccount

	cmd := `SELECT
				user_id,
				username,
				date_created,
				consumer_key,
				consumer_secret,
				access_token,
				access_token_secret
			FROM twitter_accounts
			WHERE id = $1`

	err := db.QueryRow(cmd, id).
		Scan(&account.UserID,
			&account.Username,
			&account.DateCreated,
			&account.ConsumerKey,
			&account.ConsumerSecret,
			&account.AccessToken,
			&account.AccessTokenSecret)
	if err == sql.ErrNoRows {
		return account, ErrEntityNotFound
	} else if err != nil {
		return account, err
	}

	account.ID = id
	return account, nil
}

// TwitterAccountsAll returns all TwitterAccount records from the database
var TwitterAccountsAll = func() ([]TwitterAccount, error) {
	var accounts []TwitterAccount

	cmd := `SELECT
				id,
				user_id,
				username,
				date_created,
				consumer_key,
				consumer_secret,
				access_token,
				access_token_secret
			FROM twitter_accounts`

	rows, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		account := TwitterAccount{}
		err = rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Username,
			&account.DateCreated,
			&account.ConsumerKey,
			&account.ConsumerSecret,
			&account.AccessToken,
			&account.AccessTokenSecret)

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
			tweet := Tweet{Account: account}
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
