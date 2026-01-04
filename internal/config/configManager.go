package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("cannot unmarshal config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("cannot unmarshal config file: %v", err)
	}

	return config, nil
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUser = userName
	if err := write(*c); err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFilePath := homeDir + "/" + configFileName

	return configFilePath, nil
}

/**
* 0644
* │││
* ││└─ Others (everyone else):  4 = read only
* │└── Group:                   4 = read only
* └─── Owner:                   6 = read + write
 */
func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("cannot set current user: %v", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal config: %v", err)
	}

	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		return fmt.Errorf("cannot write config: %v", err)
	}

	return nil
}
