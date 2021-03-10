package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

const ConfigENV = "API_CONFIG"

var (
	ErrConfigNotProvided = errors.New("config file is not provided")
)

type Config struct {
	Log struct {
		ShowTimestamp bool `json:"showTimestamp"`
	} `json:"log"`

	Database struct {
		MySQL struct {
			DSN          string   `json:"dsn"`
			QueryTimeout Duration `json:"queryTimeout"`
		} `json:"mysql"`
	} `json:"database"`

	HTTP struct {
		Port string `json:"port"`
		Auth struct {
			BearerToken string `json:"bearerToken"`
		} `json:"auth"`
		EnablePprof bool `json:"enablePprof"`
	} `json:"http"`
}

func getConfigPath() (string, error) {
	path := os.Getenv(ConfigENV)
	path = strings.TrimSpace(path)

	if len(path) == 0 {
		return "", ErrConfigNotProvided
	}

	return path, nil
}

func parseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = json.NewDecoder(file).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
