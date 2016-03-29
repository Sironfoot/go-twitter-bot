package db

import "time"

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
