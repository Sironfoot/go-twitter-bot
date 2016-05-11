package config_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/sironfoot/go-twitter-bot/lib/config"
)

type configuration struct {
	StringValue string  `json:"stringValue"`
	IntValue    int     `json:"intValue"`
	FloatValue  float64 `json:"floatValue"`
	BoolValue   bool    `json:"boolValue"`

	SliceValueStrings []string  `json:"sliceValueStrings"`
	SliceValueFloats  []float64 `json:"sliceValueFloats"`
	SliceValueInts    []int     `json:"sliceValueInts"`
	SliceValueBools   []bool    `json:"sliceValueBools"`

	ObjectValue       subConfiguration `json:"objectValue"`
	SliceValueObjects []slicedConfig   `json:"sliceValueObjects"`
}

type subConfiguration struct {
	StringValue string              `json:"stringValue"`
	IntValue    int                 `json:"intValue"`
	ObjectValue subSubConfiguration `json:"objectValue"`
}

type subSubConfiguration struct {
	StringValue string `json:"stringValue"`
	IntValue    int    `json:"intValue"`
}

type slicedConfig struct {
	StringValue string `json:"stringValue"`
	IntValue    int    `json:"intValue"`
}

var expectedConfig = configuration{
	StringValue: "Hello world",
	IntValue:    123,
	FloatValue:  123.45,
	BoolValue:   true,

	SliceValueStrings: []string{"string1", "string2", "string3"},
	SliceValueFloats:  []float64{1.2, 2.3, 3.4},
	SliceValueInts:    []int{1, 2, 3},
	SliceValueBools:   []bool{true, false, true},

	ObjectValue: subConfiguration{
		StringValue: "Hello world",
		IntValue:    123,

		ObjectValue: subSubConfiguration{
			StringValue: "Hello world",
			IntValue:    123,
		},
	},

	SliceValueObjects: []slicedConfig{
		slicedConfig{
			StringValue: "Hello world",
			IntValue:    123,
		},
		slicedConfig{
			StringValue: "Hello world",
			IntValue:    123,
		},
		slicedConfig{
			StringValue: "Hello world",
			IntValue:    123,
		},
	},
}

func TestLoad(t *testing.T) {
	// arrange
	var actualConfig configuration

	// act
	err := config.Load("config.json", "dev", &actualConfig)

	// assert
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedConfig, actualConfig) {
		t.Errorf("expected and actual config are different.\nExpected:\n%v\n\nActual:\n%v", expectedConfig, actualConfig)
	}
}

func TestLoadWithAltConfig(t *testing.T) {
	// arrange
	var actualConfig configuration

	altConfigString := `
    {
        "stringValue": "Hello world 2",
        "intValue": 456,
        "floatValue": 456.78,
        "boolValue": false,

        "sliceValueStrings": [ "string4", "string5", "string6", "string7" ],
		"sliceValueFloats": [ 4.5, 5.6, 7.8, 8.9 ],
		"sliceValueInts": [ 4, 5 ],
		"sliceValueBools": [ true ],

        "objectValue": {
            "stringValue": "Hello world 2",
            "intValue": 456,

            "objectValue": {
                "stringValue": "Hello world 2",
                "intValue": 456
            }
        },

        "sliceValueObjects": [
            {
                "stringValue": "Hello world 2",
                "intValue": 456
            },
            {
                "stringValue": "Hello world 2",
                "intValue": 456
            },
            {
                "stringValue": "Hello world 2",
                "intValue": 456
            },
            {
                "stringValue": "Hello world 2",
                "intValue": 456
            }
        ]
    }`

	var expectedAltConfig configuration
	err := json.Unmarshal([]byte(altConfigString), &expectedAltConfig)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile("config.test.json", []byte(altConfigString), 0644)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = os.Remove("config.test.json")
		if err != nil {
			t.Fatal(err)
		}
	}()

	// act
	err = config.Load("config.json", "test", &actualConfig)

	// assert
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedAltConfig, actualConfig) {
		t.Errorf("expected and actual config are different.\nExpected:\n%v\n\nActual:\n%v", expectedAltConfig, actualConfig)
	}
}
