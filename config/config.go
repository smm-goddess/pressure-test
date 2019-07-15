package config

import (
	"encoding/json"
	"io/ioutil"
)

func LoadConfig(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	dataJson := []byte(data)

	if err = json.Unmarshal(dataJson, v); err != nil {
		return err
	}
	return nil
}
