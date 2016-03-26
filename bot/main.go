package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type configuration struct {
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
var addr = flag.String("addr", "localhost:7000", "Address to run server on")

var (
	ticker  *time.Ticker
	stop    = make(chan bool)
	running = false
)

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			flag.PrintDefaults()
			log.Fatalln(fatalErr)
		}
	}()

	flag.Parse()

	var config configuration

	file, err := os.Open(*configFile)
	if err != nil {
		fatalErr = fmt.Errorf("can't open %s: %s", *configFile, err)
		return
	}

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		fatalErr = fmt.Errorf("can't decode %s: %s", *configFile, err)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/status", func(res http.ResponseWriter, req *http.Request) {
		if running {
			fmt.Fprint(res, "Running\n")
		} else {
			fmt.Fprint(res, "Paused\n")
		}
	})

	mux.HandleFunc("/start", func(res http.ResponseWriter, req *http.Request) {
		if !running {
			ticker = time.NewTicker(time.Second * 10)
			running = true

			go func() {
				for {
					select {
					case <-ticker.C:
						fatalErr := postNextTweet(config)
						if fatalErr != nil {
							panic(fatalErr)
						}
					case <-stop:
						return
					}
				}
			}()
		}

		fmt.Fprint(res, "Started\n")
	})

	mux.HandleFunc("/stop", func(res http.ResponseWriter, req *http.Request) {
		if running {
			ticker.Stop()
			stop <- true
			running = false
		}

		fmt.Fprint(res, "Stopped\n")
	})

	server := http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	log.Printf("Go Twitter Bot Server is running on %s...\n\n", *addr)
	server.ListenAndServe()
}

func postNextTweet(config configuration) error {
	tweets, err := LoadTweets(*dataFile)
	if err != nil {
		return fmt.Errorf("problem loading tweets: %s", err)
	}

	nextTweets := getNextTweets(tweets)

	for _, tweet := range nextTweets {
		log.Printf("Tweeting: %s\n\n", tweet.Text)

		err = postTweet(config.TwitterAuth, tweet.Text)
		if err != nil {
			return err
		}

		tweet.IsPosted = true
	}

	err = SaveTweets(tweets, *dataFile)
	if err != nil {
		return fmt.Errorf("problem saving tweets: %s", err)
	}

	return nil
}

func getNextTweets(tweets []Tweet) []*Tweet {
	now := time.Now().UTC()
	var nextTweets []*Tweet

	for i := range tweets {
		tweet := &tweets[i]

		if !tweet.IsPosted && now.After(tweet.PostOn) {
			nextTweets = append(nextTweets, tweet)
		}
	}

	return nextTweets
}
