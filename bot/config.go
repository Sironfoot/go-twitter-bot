package main

import (
	"encoding/json"
	"fmt"
	"os"
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

func loadConfig(path string) (configuration, error) {
	var config configuration

	file, err := os.Open(path)
	if err != nil {
		return config, fmt.Errorf("can't open %s: %s", *configFile, err)
	}

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return config, fmt.Errorf("can't decode %s: %s", *configFile, err)
	}

	return config, nil
}
