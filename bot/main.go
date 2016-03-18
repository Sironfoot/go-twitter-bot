package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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

	for i, tweet := range tweets {
		fmt.Printf("%v - %s\n", i+1, tweet)
	}
}
