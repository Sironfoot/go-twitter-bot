package db

import (
	"time"

	"github.com/sironfoot/go-twitter-bot/lib/sqlboiler"
)

// Tweet maps to tweets table
type Tweet struct {
	ID          string    `db:"id"`
	AccountID   string    `db:"twitter_account_id"`
	Tweet       string    `db:"tweet"`
	PostOn      time.Time `db:"post_on"`
	IsPosted    bool      `db:"is_posted"`
	DateCreated time.Time `db:"date_created"`
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
	return sqlboiler.EntitySave(tweet, dbx)
}

// Save saves the Tweet struct to the database.
func (tweet *Tweet) Save() error {
	return TweetSave(tweet)
}

// TweetDelete deletes the Tweet from the database
var TweetDelete = func(tweet *Tweet) error {
	return sqlboiler.EntityDelete(tweet, dbx)
}

// Delete deletes the Tweet from the database
func (tweet *Tweet) Delete() error {
	return TweetDelete(tweet)
}
