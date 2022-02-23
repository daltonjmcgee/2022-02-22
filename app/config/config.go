package config

import (
	"encoding/json"
	"goserve/helpers"
)

func ReturnConfig() map[string]string {
	var config map[string]string
	configJson, _ := helpers.LoadFile("./config.dev.json")
	json.Unmarshal([]byte(configJson), &config)
	return config
}
