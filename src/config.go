package main

import (
	"encoding/json"
)

func returnConfig() (map[string]string, error) {
	var config map[string]string
	configJson, _ := loadFile("../config.dev.json")
	err := json.Unmarshal([]byte(configJson), &config)
	return config, err
}
