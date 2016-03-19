package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Tweet keeps record of a tweet and whether or not it has been posted to Twitter
type Tweet struct {
	Text     string `json:"text"`
	IsPosted bool   `json:"isPosted"`
}

// LoadTweets loads Tweet structs from a json data file
func LoadTweets(dataFile string) ([]Tweet, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var tweets []Tweet
	err = json.Unmarshal(data, &tweets)

	return tweets, err
}

// SaveTweets saved an array of Tweet structs to a json data file
func SaveTweets(tweets []Tweet, dataFile string) error {
	tweetData, err := json.MarshalIndent(&tweets, "", "\t")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(dataFile, tweetData, os.ModePerm); err != nil {
		return err
	}

	return nil
}
