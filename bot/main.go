package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type config struct {
	TwitterAuth twitterAuth `json:"twitterAuth"`
}

type twitterAuth struct {
	ConsumerKey       string `json:"consumerKey"`
	ConsumerSecret    string `json:"consumerSecret"`
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
}

var configFile = flag.String("config", "config.json", "path to config file")
var dataFile = flag.String("data", "tweets.json", "path to json file containing tweets")

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			flag.PrintDefaults()
			log.Fatalln(fatalErr)
		}
	}()

	flag.Parse()

	var config config

	file, err := os.Open(*configFile)
	if err != nil {
		fatalErr = fmt.Errorf("can't open %s: %s", *configFile, err)
		return
	}

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		fatalErr = fmt.Errorf("can't decode %s: %s", *configFile, err)
		return
	}

	fmt.Println("Go Twitter Bot is running...")

	ticker := time.NewTicker(time.Second * 60)
	for range ticker.C {
		fatalErr := postNextTweet(config)
		if fatalErr != nil {
			return
		}
	}
}

func postNextTweet(config config) error {
	tweets, err := LoadTweets(*dataFile)
	if err != nil {
		return fmt.Errorf("problem loading tweets: %s", err)
	}

	tweet, err := getNextTweet(tweets)
	if err == errNoMoreTweetsToPost {
		fmt.Println("skip")
		return nil
	}

	fmt.Printf("Tweeting: %s\n\n", tweet.Text)

	err = postTweet(config.TwitterAuth, tweet.Text)
	if err != nil {
		return err
	}

	tweet.IsPosted = true

	err = SaveTweets(tweets, *dataFile)
	if err != nil {
		return fmt.Errorf("problem saving tweets: %s", err)
	}

	return nil
}

var errNoMoreTweetsToPost = errors.New("No more tweets left to be posted")

func getNextTweet(tweets []Tweet) (*Tweet, error) {
	now := time.Now().UTC()

	for i := range tweets {
		tweet := &tweets[i]

		if !tweet.IsPosted && now.After(tweet.PostOn) {
			return tweet, nil
		}
	}

	return nil, errNoMoreTweetsToPost
}
