package main

import (
	"fmt"
	"net/url"

	"github.com/mrjones/oauth"
)

func postTweet(auth twitterAuth, tweet string) error {
	consumer := oauth.NewConsumer(auth.ConsumerKey, auth.ConsumerSecret, oauth.ServiceProvider{})

	accessToken := oauth.AccessToken{
		Token:  auth.AccessToken,
		Secret: auth.AccessTokenSecret,
	}

	client, err := consumer.MakeHttpClient(&accessToken)
	if err != nil {
		return fmt.Errorf("error posting to twitter: %s", err)
	}

	_, err = client.PostForm("https://api.twitter.com/1.1/statuses/update.json", url.Values{"status": []string{tweet}})
	if err != nil {
		return fmt.Errorf("error posting to twitter: %s", err)
	}

	return nil
}
