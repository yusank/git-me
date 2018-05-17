package utils

import (
	"fmt"

	"github.com/bitly/go-simplejson"
)

// LoadJSON unmarshal json to go interface
func LoadJSON(url string, header map[string]string) (*simplejson.Json, error) {
	body, err := HttpGetByte(url, header)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))
	return simplejson.NewJson(body)
	// var f interface{}
	// if err := json.Unmarshal(data, &f); err != nil {
	// 	return nil, err
	// }

	// m := f.(map[string]interface{})
	// return m, nil
}
