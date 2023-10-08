package config

import (
	"nations/redis"

	"encoding/json"
	"os"
)

type Config struct {
	Channel string `json:"channel"`
	Guild string `json:"guild"`
}

func Get(key string) redis.Json {
	config := GetAll()
	return config[key].(redis.Json)
}

func GetStr(key string) string {
	config := GetAll()
	return config[key].(string)
}

func Save(key string, value string) {
	config := GetAll()
	config[key] = value
	dat, err := json.Marshal(config)
	if err != nil {
		panic("error saving config file (json)")
	}
	err = os.WriteFile("config.json", dat, 0644)
	if err != nil {
		panic("error saving config file (writing)")
	}
}

func GetAll() redis.Json {
	dat, err := os.ReadFile("config.json")
	if err != nil {
		panic("error reading config file")
	}
	var config redis.Json
	err = json.Unmarshal(dat, &config)
	if err != nil {
		panic("error reading config file")
	}
	return config
}
