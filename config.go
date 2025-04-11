package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TagName string `yaml:"last_downloaded_tag"`
}

func createConfigIdentifier(repo string, filename string, destination string) string {
	var cleanRepo string
	for i := range len(repo) {
		if repo[i] == '/' {
			cleanRepo = repo[:i] + repo[i+1:]
			break
		}
	}

	// shorten asset name
	if len(filename) > 20 {
		filename = filename[:20]
	}

	// use 3rd to last directory name (this is so users in prism launcher will have the instance name included in identifer)
	instName := filepath.Base(destination)
	pathParts := strings.Split(filepath.Clean(destination), string(filepath.Separator))

	if len(pathParts) >= 3 {
		instName = pathParts[len(pathParts)-3]
	}

	// create config file identifier
	configIdentifier := cleanRepo + filename + instName
	return configIdentifier
}

func getConfigPath(repo string, filename string, destination string) string {
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

	return filepath.Join(configDir, createConfigIdentifier(repo, filename, destination)+".yaml")
}

func writeConfig(config Config, repo string, filename string, destination string) error {
	configPath := getConfigPath(repo, filename, destination)

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

func loadConfig(repo string, filename string, destination string) Config {
	configPath := getConfigPath(repo, filename, destination)

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
		err := writeConfig(config, repo, filename, destination)
		if err != nil {
			fmt.Printf("Error writing config: %v\n", err)
			os.Exit(1)
		}
	}
	return config
}
