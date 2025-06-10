package records

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Record struct {
	TagName        string `yaml:"last_downloaded_tag"`
	RepositoryName string `'yaml:"repository_name"`
	FileName       string `yaml:"file_name"`
	AuthorName     string `yaml:"author_name"`
	DownloadPath   string `'yaml:"download_path"`
	UUID           string `yaml:"uuid"`
}

func getRecordDir() string {
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

	return recordDir
}

func getUUID(record Record) string {
	// generate UUID for record file
	var combinedString string = record.RepositoryName + record.FileName + record.DownloadPath

	// use combined string to generate a Version 3 / MD5 UUID
	var uuidString string = uuid.NewMD5(uuid.Nil, []byte(combinedString)).String()

	return uuidString
}

func getPath(record Record) (string, error) {
	var recordDir string = getRecordDir()

	err := os.MkdirAll(recordDir, 0755)
	if err != nil {
		return "", err
	}

	var uuidString string = getUUID(record)
	return filepath.Join(recordDir, uuidString+".yaml"), nil
}

// GetRecords returns an array of all records stored in the records/config directory
func GetRecords() ([]Record, error) {
	var records []Record

	recordDir := getRecordDir()
	files, err := os.ReadDir(recordDir)
	if err != nil {
		return nil, fmt.Errorf("error reading record directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue // skip directories they are not records
		}

		recordPath := filepath.Join(recordDir, file.Name())
		data, err := os.ReadFile(recordPath)
		if err != nil {
			return nil, fmt.Errorf("error reading record file %s: %v", file.Name(), err)
		}

		var record Record
		if err := yaml.Unmarshal(data, &record); err != nil {
			return nil, fmt.Errorf("error unmarshalling record file %s: %v", file.Name(), err)
		}

		records = append(records, record)
	}

	return records, nil
}

// Write saves a record file to the records/config directory
func Write(record Record) error {
	recordPath, err := getPath(record)
	if err != nil {
		return err
	}

	// save the UUID to the record
	record.UUID = getUUID(record)

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

// Load retrieves a record from the records/config directory,
// and begins the creation of a new record if it does not exist.
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

// Remove deletes specifed record and its respective asset file
func Remove(r string, rs []Record) error {
	// convert to int
	rInt, err := strconv.Atoi(r)
	if err != nil {
		return fmt.Errorf("error converting to int: %v", err)
	}

	// check provided number is within range of records
	if rInt > len(rs) || rInt < 1 {
		return fmt.Errorf("value is out of range (%d), record does not exist", rInt)
	}
	record := rs[rInt-1]

	// check if asset file exists and remove if it does
	path := record.DownloadPath + record.FileName
	_, err = os.Stat(path)
	notExists := os.IsNotExist(err)
	if err != nil && !notExists {
		return fmt.Errorf("error checking stat for file %s: %v", path, err)
	}

	if !notExists {
		err := os.Remove(path)
		if err != nil {
			return fmt.Errorf("error deleting file %s: %v", path, err)
		}
	}

	// delete recordfile
	path, err = getPath(record)
	if err != nil {
		return fmt.Errorf("error getting record path: %s: %v", path, err)
	}

	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("error deleting record file: %s: %v", path, err)
	}

	return nil
}
