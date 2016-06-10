package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/sironfoot/go-twitter-bot/lib/config"
)

func BenchmarkLoad(b *testing.B) {
	err := copyFile("complex.test.json", "complex.json")
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		err = os.Remove("complex.test.json")
		if err != nil {
			b.Fatal(err)
		}
	}()

	for n := 0; n < b.N; n++ {
		var complexData complex
		err := config.Load("complex.json", "test", &complexData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLoadWithCaching(b *testing.B) {
	defaultDuration := config.ReloadPollingInterval
	defer func() {
		config.SetReloadPollingInterval(defaultDuration)
	}()

	// arrange
	config.SetReloadPollingInterval(time.Duration(time.Second * 1))

	err := copyFile("complex.test.json", "complex.json")
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		err = os.Remove("complex.test.json")
		if err != nil {
			b.Fatal(err)
		}
	}()

	for n := 0; n < b.N; n++ {
		var complexData complex
		err := config.LoadWithCaching("complex.json", "test", &complexData)
		if err != nil {
			b.Fatal(err)
		}
	}
}
