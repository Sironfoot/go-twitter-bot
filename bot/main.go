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
var frequency = flag.Int64("freq", 1440, "Frequency to post each tweet (in minutes)")

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

	ticker := time.NewTicker(time.Minute * time.Duration(*frequency))

	for range ticker.C {
		tweets, err := LoadTweets(*dataFile)
		if err != nil {
			fatalErr = fmt.Errorf("problem loading tweets: %s", err)
			return
		}

		tweet, err := getNextTweet(tweets)
		if err == errNoMoreTweetsToPost {
			fmt.Println("That's all the tweets")
			continue
		}

		fmt.Printf("Tweeting: %s\n\n", tweet.Text)

		err = postTweet(config.TwitterAuth, tweet.Text)
		if err != nil {
			fatalErr = err
			return
		}

		tweet.IsPosted = true

		err = SaveTweets(tweets, *dataFile)
		if err != nil {
			fatalErr = fmt.Errorf("problem saving tweets: %s", err)
			return
		}
	}
}

var errNoMoreTweetsToPost = errors.New("No more tweets left to be posted")

func getNextTweet(tweets []Tweet) (*Tweet, error) {
	for i := range tweets {
		tweet := &tweets[i]
		if !tweet.IsPosted {
			return tweet, nil
		}
	}

	return nil, errNoMoreTweetsToPost
}
