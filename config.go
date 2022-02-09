package main

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	Test struct {
		Mode bool       `json:"mode"`
		Keys ClientKeys `json:"keys"`
	}
	Keys   ClientKeys `json:"keys"`
	Symbol string     `json:"symbol"`
}

type ClientKeys struct {
	API    string `json:"API"`
	Secret string `json:"Secret"`
}

func ReadConfig() error {
	if data, err := ioutil.ReadFile("config.json"); err == nil {
		json.Unmarshal(data, &Config)
		return nil
	} else {
		return err
	}
}
