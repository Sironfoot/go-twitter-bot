package db

import (
	"fmt"
	"time"

	"github.com/sironfoot/go-twitter-bot/lib/sqlboiler"
)

// Tweet maps to tweets table
type Tweet struct {
	ID          string
	Account     *TwitterAccount
	Tweet       string
	PostOn      time.Time
	IsPosted    bool
	DateCreated time.Time
}

// IsTransient determines if Tweet record has been saved to the database,
// true means Tweet struct has NOT been saved, false means it has.
func (tweet *Tweet) IsTransient() bool {
	return len(tweet.ID) == 0
}

// MetaData returns meta data information about the Tweet entity
func (tweet *Tweet) MetaData() sqlboiler.EntityMetaData {
	return sqlboiler.EntityMetaData{
		TableName:      "tweets",
		PrimaryKeyName: "id",
	}
}

// TweetSave saves the Tweet struct to the database.
var TweetSave = func(tweet *Tweet) error {
	if tweet.IsTransient() {
		if tweet.Account == nil {
			return fmt.Errorf("Parent TwitterAccount entity (Account field) must be set")
		}

		cmd := `INSERT INTO tweets(twitter_acount_id, tweet, post_on, is_posted, date_created)
		        VALUES($1, $2, $3, $4, $5)
                RETURNING id`

		statement, err := db.Prepare(cmd)
		if err != nil {
			return err
		}
		defer statement.Close()

		err = statement.
			QueryRow(tweet.Account.ID, tweet.Tweet, tweet.PostOn, tweet.IsPosted, tweet.DateCreated).
			Scan(&tweet.ID)
		if err != nil {
			return err
		}
	} else {
		cmd := `UPDATE tweets
		        SET tweet = $2, post_on = $3, is_posted = $4, date_created = $5
			    WHERE id = $1`

		_, err := db.Exec(cmd, tweet.ID, tweet.Tweet, tweet.PostOn, tweet.IsPosted, tweet.DateCreated)
		if err != nil {
			return err
		}
	}

	return nil
}

// Save saves the Tweet struct to the database.
func (tweet *Tweet) Save() error {
	return TweetSave(tweet)
}

// TweetDelete deletes the Tweet from the database
var TweetDelete = func(tweet *Tweet) error {
	return sqlboiler.EntityDelete(tweet, db)
}

// Delete deletes the Tweet from the database
func (tweet *Tweet) Delete() error {
	return TweetDelete(tweet)
}
