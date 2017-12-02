package utils

import (
	"github.com/bitly/go-simplejson"
)

// LoadJSON unmarshal json to go interface
func LoadJSON(data []byte) (*simplejson.Json, error) {
	return simplejson.NewJson(data)
	// var f interface{}
	// if err := json.Unmarshal(data, &f); err != nil {
	// 	return nil, err
	// }

	// m := f.(map[string]interface{})
	// return m, nil
}
