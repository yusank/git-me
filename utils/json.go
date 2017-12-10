package utils

import (
	"github.com/bitly/go-simplejson"
	"fmt"
)

// LoadJSON unmarshal json to go interface
func LoadJSON(url string, header map[string]string) (*simplejson.Json, error) {
	body, err := GetContent(url, header)
	if err != nil {
		return nil, err
	}

	fmt.Println(len(body))
	fmt.Println(string(body))
	return simplejson.NewJson(body)
	// var f interface{}
	// if err := json.Unmarshal(data, &f); err != nil {
	// 	return nil, err
	// }

	// m := f.(map[string]interface{})
	// return m, nil
}
