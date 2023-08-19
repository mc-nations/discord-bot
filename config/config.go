package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Channel string `json:"channel"`
}

func Get(key string) string {
	config := GetAll()
	switch key {
	case "channel":
		return config.Channel
	}
	return ""
}

func Save(key string, value string) {
	config := GetAll()
	switch key {
	case "channel":
		config.Channel = value
	}
	dat, err := json.Marshal(config)
	if err != nil {
		panic("error saving config file (json)")
	}
	err = os.WriteFile("config.json", dat, 0644)
	if err != nil {
		panic("error saving config file (writing)")
	}
}

func GetAll() Config {
	dat, err := os.ReadFile("config.json")
	if err != nil {
		panic("error reading config file")
	}
	var config Config
	err = json.Unmarshal(dat, &config)
	if err != nil {
		panic("error reading config file")
	}
	return config
}
