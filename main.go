package main

import (
	"encoding/json"
	"io"

	"fmt"
	"net/http"
	"os"
)

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
}

func main() {
	usageString := "Usage: githubpackdownloader <repo> <filename> [<repo> <filename>] <destination>"

	// check if we have repo+filename+destination at minimum
	if (len(os.Args) - 1) < 3 {
		fmt.Println(usageString)
		os.Exit(1)
	}

	// check if odd number of args (missing filename for a repo)
	if (len(os.Args)-2)%2 != 0 {
		fmt.Println(usageString)
		os.Exit(1)
	}

	destination := os.Args[len(os.Args)-1]

	if destination[len(destination)-1:] != "/" {
		destination += "/"
	}

	// loop through all given repo + filename pairs
	for i := 1; i < len(os.Args)-1; i += 2 {
		repo := os.Args[i]
		filename := os.Args[i+1]

		// get latest release, in separate function to isolate defer calls
		release, err := getRelease(repo)
		if err != nil {
			fmt.Printf("Error getting release: %v\n", err)
			os.Exit(1)
		}

		// check if already exists
		record, err := loadRecord(repo, filename, destination)
		if err != nil {
			fmt.Printf("Error loading record: %v\n", err)
			os.Exit(1)
		}
		denyDownload := false
		if _, err := os.Stat(destination + filename); err == nil {
			// check version match
			if record.TagName == release.TagName {
				denyDownload = true
			}
		}

		if !denyDownload {
			// acquire asset
			var file Asset
			for _, asset := range release.Assets {
				if asset.Name == filename {
					file = asset
					break
				}
			}

			// download file
			err := downloadFile(destination, filename, file.DownloadURL)
			if err != nil {
				fmt.Printf("Error downloading file: %v\n", err)
				os.Exit(1)
			}

			// update yaml record with new info
			record.TagName = release.TagName
			err = writeRecord(record, repo, filename, destination)
			if err != nil {
				fmt.Printf("Error writing config file: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

// acquire the latest release from the given GitHub repo
func getRelease(repo string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("response returned with error: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
			os.Exit(1)
		}
	}(response.Body)

	var release Release
	if err := json.NewDecoder(response.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return &release, nil
}

func downloadFile(destination string, filename string, url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
			os.Exit(1)
		}
	}(response.Body)

	// create file
	output, err := os.Create(destination + filename)
	if err != nil {
		return err
	}
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {
			fmt.Printf("Error closing output file: %v\n", err)
			os.Exit(1)
		}
	}(output)

	// write information to file
	_, err = io.Copy(output, response.Body)
	return err
}
