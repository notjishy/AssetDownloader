package records

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Record struct {
	TagName        string `yaml:"last_downloaded_tag"`
	RepositoryName string `'yaml:"repository_name"`
	FileName       string `yaml:"file_name"`
	AuthorName     string `yaml:"author_name"`
	DownloadPath   string `'yaml:"download_path"`
}

func getPath(record Record) (string, error) {
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

	// generate UUID for record file
	var combinedString string = record.RepositoryName + record.FileName + record.DownloadPath

	// use combined string to generate a Version 3 / MD5 UUID
	var uuidString string = uuid.NewMD5(uuid.Nil, []byte(combinedString)).String()
	return filepath.Join(recordDir, uuidString+".yaml"), nil
}

func Write(record Record) error {
	recordPath, err := getPath(record)
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

func Load(repo string, filename string, destination string) (Record, error) {
	var record Record = Record{}

	record.RepositoryName = repo
	record.FileName = filename
	record.DownloadPath = destination

	recordPath, err := getPath(record)
	if err != nil {
		return Record{}, err
	}

	// check if record file exists
	data, err := os.ReadFile(recordPath)
	if err != nil {
		if os.IsNotExist(err) {
			// file not does not exist create new record
			if err := Write(record); err != nil {
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
