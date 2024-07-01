package internals

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Path        string
	DataSources string `json:"datasources"`
}

var appConfig *Config

func LoadConfig(path string) *Config {
	var config *Config

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal(err)
	}

	config.Path = path
	appConfig = config

	return config
}

func GetConfig() *Config {
	if appConfig == nil {
		panic("config not initialized")
	}

	return appConfig
}
