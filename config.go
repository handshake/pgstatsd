package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	ConnectionURL string
}

func ReadConfig(cfgPath string) Configuration {
	var config Configuration

	file, err := os.Open(cfgPath)
	defer file.Close()
	if err != nil {
		log.Fatalf("couldn't read config from %s", cfgPath)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("couldn't read config from conf.json")
	}
	return config
}
