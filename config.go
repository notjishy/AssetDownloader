package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

func getConfigPath() string {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		// %APPDATA%
		configDir = filepath.Join(os.Getenv("APPDATA"), "GithubPackDownloader")
	case "darwin": // macOS
		// ~/Library/Application Support/
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, "Library", "Application Support", "GithubPackDownloader")
	default: // Linux
		// ~/.config/
		configDir = filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "GithubPackDownloader")
		if configDir == filepath.Join("", "GithubPackDownloader") {
			// XDG_CONFIG_HOME not set, fallback to ~/.config
			homeDir, _ := os.UserHomeDir()
			configDir = filepath.Join(homeDir, ".config", "GithubPackDownloader")
		}
	}

	os.MkdirAll(configDir, 0755)

	return filepath.Join(configDir, ".ghpdconfig.yaml")
}

func writeConfig(config Config, configPath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}
	return err
}

func loadConfig(configPath string) Config {
	config := Config{}
	// check if config file exists
	if _, err := os.Stat(configPath); err == nil {
		// does exists
		data, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Error reading YAML file: %v\n", err)
			os.Exit(1)
		}

		err = yaml.Unmarshal(data, &config)
		if err != nil {
			fmt.Printf("Error unmarshaling YAML: %v\n", err)
			os.Exit(1)
		}
	} else {
		// does not exist
		config.TagName = ""
		err := writeConfig(config, configPath)
		if err != nil {
			fmt.Printf("Error writing config: %v\n", err)
			os.Exit(1)
		}
	}
	return config
}
