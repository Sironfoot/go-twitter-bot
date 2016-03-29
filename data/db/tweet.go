package db

import (
	"fmt"
	"time"
)

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

// IsTransient determines if Tweet record has been saved to the database,
// true means Tweet struct has NOT been saved, false means it has.
func (tweet *Tweet) IsTransient() bool {
	return len(tweet.id) == 0
}

// Save saves the Tweet struct to the database.
func (tweet *Tweet) Save() error {
	if tweet.IsTransient() {
		if tweet.Account == nil {
			return fmt.Errorf("Parent TwitterAccount entity must be set")
		}

		sql := "INSERT INTO tweets(twitter_acount_id, tweet, post_on, is_posted, date_created) " +
			"VALUES($1, $2, $3, $4, $5) RETURNING id"

		statement, err := db.Prepare(sql)
		if err != nil {
			return err
		}
		defer statement.Close()

		err = statement.
			QueryRow(tweet.Account.id, tweet.Tweet, tweet.PostOn, tweet.IsPosted, tweet.DateCreated).
			Scan(&tweet.id)
		if err != nil {
			return err
		}
	} else {
		_, err := db.Exec("UPDATE tweets SET tweet = $2, post_on = $3, is_posted = $4, date_created = $5 WHERE id = $1",
			tweet.id, tweet.Tweet, tweet.PostOn, tweet.IsPosted, tweet.DateCreated)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes the Tweet from the database
func (tweet *Tweet) Delete() error {
	_, err := db.Exec("DELETE FROM tweets WHERE id = $1", tweet.id)
	return err
}
