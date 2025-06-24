package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

/**
 * Config path:
 * - Linux: ~/.config/passenger-go/config.json
 * - macOS: ~/Library/Application Support/passenger-go/config.json
 * - Windows: %APPDATA%/passenger-go/config.json
 */

type Config struct {
	ServerURL string `json:"server_url,omitempty"`
}

func getConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "passenger-go", "config.json"), nil
}

func LoadConfig() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	config := &Config{}

	reader, err := os.Open(path)
	if err != nil {
		SaveConfig(config)
	}
	defer reader.Close()

	json.NewDecoder(reader).Decode(config)
	return config, nil
}

func SaveConfig(config *Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	writer, err := os.Create(path)
	if err != nil {
		return err
	}
	defer writer.Close()

	json.NewEncoder(writer).Encode(config)
	return nil
}
