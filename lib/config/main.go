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

// Load will load a configuration json file into a struct
func Load(path, mode string, configData interface{}) (err error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return ErrPrimaryConfigFileNotExist
	} else if err != nil {
		return err
	}

	err = json.NewDecoder(file).Decode(configData)
	if err != nil {
		return err
	}

	altPath := strings.Replace(path, ".json", "."+mode+".json", 1)

	altFile, err := os.Open(altPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		return nil
	}

	altConfigData := map[string]interface{}{}
	err = json.NewDecoder(altFile).Decode(&altConfigData)
	if err != nil {
		return err
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

		switch realValue := value.(type) {
		case map[string]interface{}:
			parseMap(realValue, fieldValue)
		case []interface{}:
			parseSlice(realValue, fieldValue)
		case string:
			fieldValue.SetString(realValue)
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
			fieldValue.SetBool(realValue)
		}
	}
}

func parseSlice(aSlice []interface{}, configValue reflect.Value) {
	configValue.SetLen(0)
	configValue.SetCap(0)
	newSlice := reflect.MakeSlice(reflect.SliceOf(configValue.Type().Elem()), len(aSlice), len(aSlice))
	configValue.Set(newSlice)

	for i, value := range aSlice {
		switch realItem := value.(type) {
		case map[string]interface{}:
			parseMap(realItem, configValue.Index(i))
		case []interface{}:
			parseSlice(realItem, configValue.Index(i))
		case string:
			configValue.Index(i).SetString(realItem)
		case float64:
			switch configValue.Index(i).Kind() {
			case reflect.Float32, reflect.Float64:
				configValue.Index(i).SetFloat(realItem)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				configValue.Index(i).SetInt(int64(realItem))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				configValue.Index(i).SetUint(uint64(realItem))
			}
		case bool:
			configValue.Index(i).SetBool(realItem)
		}
	}
}
