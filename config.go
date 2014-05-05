package main

import (
	"encoding/json"
	"log"
	"os"
)

type PostgresConfig struct {
	ConnectionString string `json:"connection_string"`
}

type StatsdConfig struct {
	ConnectionString string `json:"connection_string"`
	Prefix           string `json:"prefix"`
}

type Configuration struct {
	PG PostgresConfig `json:"postgres"`
	ST StatsdConfig   `json:"statsd"`
}

func ReadConfig(cfgPath string) Configuration {
	var config Configuration

	file, err := os.Open(cfgPath)
	defer file.Close()
	if err != nil {
		log.Fatalf("couldn't read config from %s [%s]", cfgPath, err.Error())
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("couldn't read config from conf.json [%s]", err.Error())
	}
	return config
}
