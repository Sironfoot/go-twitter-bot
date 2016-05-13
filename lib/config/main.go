// Package config provides utilities for loading a JSON configuration file into a struct object
// graph, with support for providing an alternative environment JSON config file
// (e.g. "dev", "staging", "uat", "live"), with values replaced using transformations. Inspired by
// the way Microsoft ASP.NET handles configuration files.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// ErrPrimaryConfigFileNotExist is returned if the primary JSON config file doesn't exist
var ErrPrimaryConfigFileNotExist = fmt.Errorf("config: primary config file does not exist")

// ErrConfigDataNotPointer is returned when the configData
// struct to pass config data into is not a pointer
var ErrConfigDataNotPointer = fmt.Errorf("config: configData argument is not a pointer")

// Load will load a configuration json file into a struct
func Load(path, environment string, configData interface{}) (err error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return ErrPrimaryConfigFileNotExist
	} else if err != nil {
		return fmt.Errorf("config: error opening primary config file: %s", err)
	}

	if reflect.TypeOf(configData).Kind() != reflect.Ptr {
		return ErrConfigDataNotPointer
	}

	err = json.NewDecoder(file).Decode(configData)
	if err != nil {
		return fmt.Errorf("config: cannot unmarshal config file: %s", err)
	}

	altPath := strings.Replace(path, ".json", "."+environment+".json", 1)

	altFile, err := os.Open(altPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("config: error opening environment config file \"%s\": %s", altPath, err)
	}

	if os.IsNotExist(err) {
		return nil
	}

	altConfigData := map[string]interface{}{}
	err = json.NewDecoder(altFile).Decode(&altConfigData)
	if err != nil {
		return fmt.Errorf("config: cannot unmarshal environment config file: %s", err)
	}

	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("config: parsing environment config file: %s", p)
		}
	}()

	configValue := reflect.ValueOf(configData).Elem()
	parseMap(altConfigData, configValue)

	return
}

func parseMap(aMap map[string]interface{}, configValue reflect.Value) {
	for key, value := range aMap {
		fieldName := ""

		for i := 0; i < configValue.NumField(); i++ {
			fieldInfo := configValue.Type().Field(i)
			jsonFieldName := strings.TrimSpace(fieldInfo.Tag.Get("json"))

			if jsonFieldName == key {
				fieldName = fieldInfo.Name
			}
		}

		fieldValue := configValue.FieldByName(fieldName)
		if fieldValue.Kind() == reflect.Invalid {
			continue
		}

		switch realValue := value.(type) {
		case map[string]interface{}:
			if fieldValue.Kind() == reflect.Struct || fieldValue.Kind() == reflect.Map {
				parseMap(realValue, fieldValue)
			}
		case []interface{}:
			if fieldValue.Kind() == reflect.Slice {
				parseSlice(realValue, fieldValue)
			}
		case string:
			if fieldValue.Kind() == reflect.String {
				fieldValue.SetString(realValue)
			}
		case float64:
			switch fieldValue.Kind() {
			case reflect.Float32, reflect.Float64:
				fieldValue.SetFloat(realValue)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue.SetInt(int64(realValue))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue.SetUint(uint64(realValue))
			}
		case bool:
			if fieldValue.Kind() == reflect.Bool {
				fieldValue.SetBool(realValue)
			}
		}
	}
}

func parseSlice(aSlice []interface{}, configValue reflect.Value) {
	newSlice := reflect.MakeSlice(configValue.Type(), len(aSlice), len(aSlice))
	configValue.Set(newSlice)

	for i, value := range aSlice {
		configItem := configValue.Index(i)

		switch realItem := value.(type) {
		case map[string]interface{}:
			if configItem.Kind() == reflect.Struct || configItem.Kind() == reflect.Map {
				parseMap(realItem, configItem)
			}
		case []interface{}:
			if configItem.Kind() == reflect.Slice {
				parseSlice(realItem, configItem)
			}
		case string:
			if configItem.Kind() == reflect.String {
				configItem.SetString(realItem)
			}
		case float64:
			switch configItem.Kind() {
			case reflect.Float32, reflect.Float64:
				configItem.SetFloat(realItem)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				configItem.SetInt(int64(realItem))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				configItem.SetUint(uint64(realItem))
			}
		case bool:
			if configItem.Kind() == reflect.Bool {
				configItem.SetBool(realItem)
			}
		}
	}
}
