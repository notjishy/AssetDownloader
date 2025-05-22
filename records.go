package main

import (
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

func getRecordPath(repo string, filename string, destination string) (string, error) {
	var recordDir string

	switch runtime.GOOS {
	case "windows":
		// %APPDATA%
		recordDir = filepath.Join(os.Getenv("APPDATA"), "AssetDownloader")
	case "darwin": // macOS
		// ~/Library/Application Support/
		homeDir, _ := os.UserHomeDir()
		recordDir = filepath.Join(homeDir, "Library", "Application Support", "AssetDownloader")
	default: // Linux
		// ~/.config/
		recordDir = filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "AssetDownloader")
		if recordDir == filepath.Join("", "AssetDownloader") {
			// XDG_CONFIG_HOME not set, fallback to ~/.config
			homeDir, _ := os.UserHomeDir()
			recordDir = filepath.Join(homeDir, ".config", "AssetDownloader")
		}
	}

	err := os.MkdirAll(recordDir, 0755)
	if err != nil {
		return "", err
	}

	return filepath.Join(recordDir, createRecordIdentifier(repo, filename, destination)+".yaml"), nil
}

func writeRecord(record Record, repo string, filename string, destination string) error {
	recordPath, err := getRecordPath(repo, filename, destination)
	if err != nil {
		return err
	}

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

func loadRecord(repo string, filename string, destination string) (Record, error) {
	recordPath, err := getRecordPath(repo, filename, destination)
	if err != nil {
		return Record{}, err
	}

	record := Record{}
	// check if record file exists
	data, err := os.ReadFile(recordPath)
	if err != nil {
		if os.IsNotExist(err) {
			// file not does not exist create new record
			record.TagName = ""
			if err := writeRecord(record, repo, filename, destination); err != nil {
				return Record{}, err
			}
			return record, nil
		}
		return Record{}, err
	}

	if err := yaml.Unmarshal(data, &record); err != nil {
		return Record{}, err
	}
	return record, nil
}
