package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error establising home directory: %v", err)
		return "", err
	}
	return home + "/" + configFileName, nil
}

func Read() (Config, error) {
	fpth, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(fpth)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Printf("Error Unmarshalling data: %v", err)
		return Config{}, err
	}
	return cfg, nil
}

func write(cfg Config) error {
	fpth, err := getConfigFilePath()
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		log.Printf("Error marshalling struct to json: %v", err)
		return err
	}

	err = os.WriteFile(fpth, jsonData, 0644)
	if err != nil {
		log.Printf("Error writing file: %v", err)
		return err
	}
	return nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)

}
