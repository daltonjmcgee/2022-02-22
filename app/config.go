package main

import (
	"encoding/json"
)

func returnConfig() map[string]string {
	var config map[string]string
	configJson, _ := loadFile("./config.dev.json")
	json.Unmarshal([]byte(configJson), &config)
	return config
}
