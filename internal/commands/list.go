package commands

import (
	"fmt"

	"jishe.wtf/assetdownloader/internal/records"
)

// ListRecords returns a string representation of all records.
// if an error occurs, it returns a non-empty error message
func ListRecords() (string, error) {
	r, err := records.GetRecords()
	if err != nil {
		return "", fmt.Errorf("error getting records: %v", err)
	}

	if len(r) == 0 {
		return "", fmt.Errorf("no records found")
	}

	var rString string
	for i, r := range r {
		rString += fmt.Sprint(i+1) + ": " + r.UUID + " - " + r.RepositoryName + " - "
		rString += r.TagName + " - " + r.FileName + " - " + r.DownloadPath + "\n\n"
	}

	return rString, nil
}
