package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"jishe.wtf/assetdownloader/internal/records"
)

var destination string
var finalArg string

var hasDownloadAllFlag bool

type Release struct {
	TagName string  `json:"tag_name"`
	Author  Author  `json:"author"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
}

type Author struct {
	Name string `json:"login"`
}

// Download downloads the latest release of the specified GitHub repository asset.
// If it already exists as the latest version, it will skip it.
func Download(args []string) error {
	// check for download all flag
	finalArg = args[len(args)-1]
	if finalArg == "-a" || finalArg == "--download-all" {
		hasDownloadAllFlag = true
	}

	// check if minimum number of arguments is provided and
	// if amount of arguments is NOT a value of the below sequence
	// 4, 7, 10, 13, 16, 19, 22, 25, ...
	if ((len(args) - 1) < 3) || ((len(args)-1)%3 != 0 && !hasDownloadAllFlag) {
		return fmt.Errorf("invalid number of arguments provided\n\n%s", getUsageString())
	}

	// only run if flag is at the end
	if hasDownloadAllFlag {
		if (len(args)-1)%2 != 0 {
			return fmt.Errorf("invalid number of arguments provided with download all flag\n\n%s", getUsageString())
		}

		fmt.Println("Flag identified...")
		destination = parseDownloadPath(args[len(args)-2])
	}

	// loop through all given repo + filename pairs
	for i := 1; i < len(args)-2; i += 2 {
		repo := args[i]
		filename := args[i+1]

		// get download destination for current asset
		if !hasDownloadAllFlag {
			destination = parseDownloadPath(args[i+2])
			i++
		}

		// get latest release, in separate function to isolate defer calls
		release, err := getRelease(repo)
		if err != nil {
			return fmt.Errorf("error getting release: %v", err)
		}

		// check if already exists
		record, err := records.Load(repo, filename, destination)
		if err != nil {
			return fmt.Errorf("error loading record: %v", err)
		}

		if _, err := os.Stat(destination + filename); err == nil {
			// check version match
			if record.TagName == release.TagName {
				continue // skip to next loop iteration
			}
		}

		// acquire asset data
		var file Asset
		for _, asset := range release.Assets {
			if asset.Name == filename {
				file = asset
				break
			}
		}

		// download the asset
		err = downloadAsset(destination, filename, file.DownloadURL)
		if err != nil {
			return fmt.Errorf("error downloading file: %v", err)
		}

		// update yaml record with new info
		record.TagName = release.TagName
		record.AuthorName = release.Author.Name

		err = records.Write(record)
		if err != nil {
			return fmt.Errorf("error writing record: %v", err)
		}
	}

	return nil
}

// returns the usage string for the download command.
func getUsageString() string {
	uString := "Usage: download <repo> <filename> [<repo> <filename>] <destination> <-a | --download-all>\n"
	uString += "download <repo> <filename> <destination> [<repo> <filename> <destination>]\n"

	return uString
}

// acquire the latest release from the given GitHub repo
func getRelease(repo string) (release *Release, err error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("response returned with error: %v", err)
	}
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			err = fmt.Errorf("error closing response body: %v", closeErr)
		}
	}(response.Body)

	if err := json.NewDecoder(response.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return release, nil
}

// setting destination variable and path formatting
func parseDownloadPath(path string) string {
	// ensure trailing "/" always present
	if path[len(path)-1:] != "/" {
		path += "/"
	}

	return path
}

func downloadAsset(destination string, filename string, url string) (err error) {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			err = fmt.Errorf("error closing response body: %v", closeErr)
		}
	}(response.Body)

	// create file
	output, err := os.Create(destination + filename)
	if err != nil {
		return err
	}
	defer func(output *os.File) {
		closeErr := output.Close()
		if closeErr != nil {
			err = fmt.Errorf("error closing output file: %v", closeErr)
		}
	}(output)

	// write information to file
	_, err = io.Copy(output, response.Body)
	return err
}
