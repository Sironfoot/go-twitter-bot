package main

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestLoadingData(t *testing.T) {
	tweets, err := LoadTweets("tweets.json")
	if err != nil {
		t.Fatalf("Failed to load tweets: %s", err)
	}

	if len(tweets) == 0 {
		t.Error("Should be at least 1 tweet")
	}

	fmt.Printf("%d tweets were returned\n", len(tweets))
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
