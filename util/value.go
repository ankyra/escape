/*
Copyright 2017 Ankyra

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func InterfaceMapToStringMap(values *map[string]interface{}, keyPrefix string) map[string]string {
	result := map[string]string{}
	if values == nil {
		return result
	}
	for key, val := range *values {
		stringVal, err := InterfaceToString(val)
		if err != nil {
			panic(fmt.Sprintf("%s (key: '%s'). This is a bug in Escape.", err.Error(), key))
		}
		result[keyPrefix+key] = stringVal
	}
	return result
}

func InterfaceToString(val interface{}) (string, error) {
	stringVal := ""
	switch val.(type) {
	case string:
		stringVal = val.(string)
	case bool:
		stringVal = "0"
		if val.(bool) {
			stringVal = "1"
		}
	case float64:
		stringVal = strconv.Itoa(int(val.(float64)))
	case int:
		stringVal = strconv.Itoa(val.(int))
	case []interface{}:
		jsonBytes, err := json.Marshal(val)
		if err != nil {
			panic(err)
		}
		stringVal = string(jsonBytes)
	default:
		return "", fmt.Errorf("Type '%T' not supported", val)
	}
	return stringVal, nil
}

func StructToMapStringInterface(val interface{}, tag string) map[string]interface{} {
	out := make(map[string]interface{})

	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("util.StructToMapStringInterface only accepts structs; got %T", v))
	}

	vType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := vType.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" {
			out[strings.Split(tagv, ",")[0]] = v.Field(i).Interface()
		}
	}
	return out
}
