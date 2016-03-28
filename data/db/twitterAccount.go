package db

import "time"

// TwitterAccount maps to twitter_accounts table
type TwitterAccount struct {
	id                string
	UserID            string
	Username          string
	DateCreated       time.Time
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// ID returns read-only Primary Key ID of TwitterAccount
func (account *TwitterAccount) ID() string {
	return account.id
}

// IsTransient determines if TwitterAccount record has been saved to the database,
// true means TwitterAccount struct has NOT been saved, false means it has.
func (account *TwitterAccount) IsTransient() bool {
	return len(account.id) == 0
}

// Save saves the TwitterAccount struct to the database.
func (account *TwitterAccount) Save() error {
	if account.IsTransient() {
		sql := "INSERT INTO twitter_accounts(user_id, username, date_created, consumer_key, consumer_secret, access_token, access_token_secret) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id"

		statement, err := db.Prepare(sql)
		if err != nil {
			return err
		}
		defer statement.Close()

		err = statement.
			QueryRow(account.UserID, account.Username, account.DateCreated, account.ConsumerKey, account.ConsumerSecret, account.AccessToken, account.AccessTokenSecret).
			Scan(&account.id)
		if err != nil {
			return err
		}
	} else {
		_, err := db.Exec("UPDATE twitter_accounts SET user_id = $2, username = $3 date_created = $4, consumer_key = $5, consumer_secret = $6, access_token = $7, access_token_secret = $8 WHERE id = $1",
			account.id, account.UserID, account.Username, account.DateCreated, account.ConsumerKey, account.ConsumerSecret, account.AccessToken, account.AccessTokenSecret)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes the TwitterAccount from the database
func (account *TwitterAccount) Delete() error {
	_, err := db.Exec("DELETE FROM twitter_accounts WHERE id = $1", account.id)
	return err
}

// Tweet maps to tweets table
type Tweet struct {
	id          string
	Account     *TwitterAccount
	Tweet       string
	PostOn      time.Time
	IsPosted    bool
	DateCreated time.Time
}

// ID returns read-only Primary Key ID of Tweet
func (tweet *Tweet) ID() string {
	return tweet.id
}

// GetTweets loads Tweets child entites for TwitterAccount
func (account *TwitterAccount) GetTweets() ([]Tweet, error) {
	var tweets []Tweet

	if !account.IsTransient() {
		rows, err := db.Query("SELECT id, tweet, post_on, is_posted, date_created FROM tweets WHERE twitter_account_id = $1", account.id)
		if err != nil {
			return tweets, err
		}

		for rows.Next() {
			tweet := Tweet{Account: account}
			err = rows.Scan(&tweet.id, &tweet.Tweet, &tweet.PostOn, &tweet.IsPosted, &tweet.DateCreated)
			if err != nil {
				return tweets, err
			}

			tweets = append(tweets, tweet)
		}

		rows.Close()
	}

	return tweets, nil
}
