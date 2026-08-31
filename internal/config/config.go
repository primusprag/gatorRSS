package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = "/.gatorconfig.json"


type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("Error opening file path: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var cfg Config
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("Error decoding json: %w", err)
	}

	return cfg, nil
}

func (cfg *Config) SetUser(userName string) error {

	cfg.CurrentUserName = userName

	err := write(*cfg)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error retrieving home path: %w", err)
	}
	return homePath + configFileName, nil
}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
