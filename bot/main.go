package main

import (
	"encoding/json"
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

func main() {
	var config config

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("can't open config.json: ", err)
	}

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal("can't decode config.json: ", err)
	}

	fmt.Println("Go Twitter Bot is running...")

	fmt.Println(config.TwitterAuth.AccessToken)
}
