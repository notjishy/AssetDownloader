package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

type Record struct {
	TagName string `yaml:"last_downloaded_tag"`
}

func createRecordIdentifier(repo string, filename string, destination string) string {
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

	// create record file identifier
	recordIdentifier := cleanRepo + filename + instName
	return recordIdentifier
}

func getRecordPath(repo string, filename string, destination string) string {
	var recordDir string

	switch runtime.GOOS {
	case "windows":
		// %APPDATA%
		recordDir = filepath.Join(os.Getenv("APPDATA"), "GithubPackDownloader")
	case "darwin": // macOS
		// ~/Library/Application Support/
		homeDir, _ := os.UserHomeDir()
		recordDir = filepath.Join(homeDir, "Library", "Application Support", "GithubPackDownloader")
	default: // Linux
		// ~/.config/
		recordDir = filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "GithubPackDownloader")
		if recordDir == filepath.Join("", "GithubPackDownloader") {
			// XDG_CONFIG_HOME not set, fallback to ~/.config
			homeDir, _ := os.UserHomeDir()
			recordDir = filepath.Join(homeDir, ".config", "GithubPackDownloader")
		}
	}

	os.MkdirAll(recordDir, 0755)

	return filepath.Join(recordDir, createRecordIdentifier(repo, filename, destination)+".yaml")
}

func writeRecord(record Record, repo string, filename string, destination string) error {
	recordPath := getRecordPath(repo, filename, destination)

	data, err := yaml.Marshal(record)
	if err != nil {
		return err
	}

	err = os.WriteFile(recordPath, data, 0644)
	if err != nil {
		return err
	}
	return err
}

func loadRecord(repo string, filename string, destination string) Record {
	recordPath := getRecordPath(repo, filename, destination)

	record := Record{}
	// check if record file exists
	if _, err := os.Stat(recordPath); err == nil {
		// does exists
		data, err := os.ReadFile(recordPath)
		if err != nil {
			fmt.Printf("Error reading YAML file: %v\n", err)
			os.Exit(1)
		}

		err = yaml.Unmarshal(data, &record)
		if err != nil {
			fmt.Printf("Error unmarshaling YAML: %v\n", err)
			os.Exit(1)
		}
	} else {
		// does not exist
		record.TagName = ""
		err := writeRecord(record, repo, filename, destination)
		if err != nil {
			fmt.Printf("Error writing config: %v\n", err)
			os.Exit(1)
		}
	}
	return record
}
