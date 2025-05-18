package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DataBaseURL string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func (cfg Config) write() error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	configFileData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	os.WriteFile(configFilePath, configFileData, 0755)

	return nil
}

func (cfg Config) SetUser(user string) error {
	cfg.CurrentUser = user
	return cfg.write()
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + "/" + configFileName, nil
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	configFileData, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(configFileData, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
