package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	SecretKey   string
	Server      string
	Zone        string
	Domain      string
	RecordTTL   int
	NsupdateCmd string
}

func (conf *Config) LoadConfig(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		panic(err)
	}
}
