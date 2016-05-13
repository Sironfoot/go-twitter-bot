package config_test

import (
	"reflect"
	"testing"

	"github.com/sironfoot/go-twitter-bot/lib/config"
)

func TestLoad_NonPointer(t *testing.T) {
	var configuration complex

	err := config.Load("complex.json", "", configuration)

	// assert
	if err != config.ErrConfigDataNotPointer {
		t.Errorf("should have returned error: %s", config.ErrConfigDataNotPointer)
	}
}

func TestLoad_NoEnvironment(t *testing.T) {
	// arrange
	var actualConfig complex

	// act
	err := config.Load("complex.json", "dev", &actualConfig)

	// assert
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedConfig, actualConfig) {
		t.Errorf("expected and actual config are different.\nExpected:\n%v\n\nActual:\n%v", expectedConfig, actualConfig)
	}
}
