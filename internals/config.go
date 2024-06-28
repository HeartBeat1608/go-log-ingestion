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
	if appConfig != nil {
		return appConfig
	}

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
