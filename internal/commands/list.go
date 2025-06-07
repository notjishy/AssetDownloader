package commands

import (
	"fmt"

	"jishe.wtf/assetdownloader/internal/records"
)

// ListRecords returns a string representation of all records.
// if an error occurs, it returns a non-empty error message
func ListRecords() (string, error) {
	records, err := records.GetRecords()
	if err != nil {
		return "", fmt.Errorf("error getting records: %v", err)
	}

	if len(records) == 0 {
		return "", fmt.Errorf("no records found")
	}

	var rString string
	for i, r := range records {
		rString += fmt.Sprint(i+1) + ": " + r.UUID + " - " + r.RepositoryName + " - "
		rString += r.TagName + " - " + r.FileName + " - " + r.DownloadPath + "\n\n"
	}

	return rString, nil
}
