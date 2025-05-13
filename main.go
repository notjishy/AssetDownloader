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
	if len(os.Args) < 4 {
		fmt.Println("Usage: githubpackdownloader <repo> <filename> <destination>")
		os.Exit(1)
	}

	repo := os.Args[1]
	filename := os.Args[2]
	destination := os.Args[3]

	if destination[len(destination)-1:] != "/" {
		destination += "/"
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Response returned with error: %v\n", err)
		os.Exit(1)
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
		fmt.Printf("Error parsing JSON: %v\n", err)
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
