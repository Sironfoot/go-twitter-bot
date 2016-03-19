package main

import (
	"os"
	"testing"
	"unicode/utf8"
)

func TestLoadTweets(t *testing.T) {
	tweets, err := LoadTweets("tweets.json")
	if err != nil {
		t.Fatalf("Failed to load tweets: %s", err)
	}

	if len(tweets) == 0 {
		t.Error("Should be at least 1 tweet")
	}
}

func TestSaveTweets(t *testing.T) {
	tweetFile := "tweets_test.json"

	tweets := []Tweet{
		Tweet{
			Text:     "Tweet 1",
			IsPosted: false,
		},
		Tweet{
			Text:     "Tweet 2",
			IsPosted: true,
		},
	}

	err := SaveTweets(tweets, tweetFile)
	if err != nil {
		t.Fatalf("Failed to save tweets: %s", err)
	}

	defer os.Remove(tweetFile)

	savedTweets, err := LoadTweets(tweetFile)
	if err != nil {
		t.Fatalf("Failed to load tweets: %s", err)
	}

	if len(savedTweets) != len(tweets) {
		t.Fatalf("Number of tweets (%d) not the same as number saved (%d)", len(tweets), len(savedTweets))
	}

	for i := range tweets {
		tweet := tweets[i]
		savedTweet := savedTweets[i]

		if tweet.Text != savedTweet.Text && tweet.IsPosted != savedTweet.IsPosted {
			t.Errorf("tweet at index: %d doesn't match saved tweet", i)
		}
	}
}

func TestTweetMaxLength(t *testing.T) {
	tweets, err := LoadTweets("tweets.json")
	if err != nil {
		t.Fatalf("Failed to load tweets: %s", err)
	}

	for i, tweet := range tweets {
		charCount := utf8.RuneCountInString(tweet.Text)
		if charCount > 140 {
			t.Errorf("Tweet at index: %d was %d characters. Tweet was:\n\n%s\n\n\n", i, charCount, tweet.Text)
		}
	}
}
