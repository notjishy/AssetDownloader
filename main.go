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
	defer response.Body.Close()

	var release Release
	if err := json.NewDecoder(response.Body).Decode(&release); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// check if already exists
	config := loadConfig()
	denyDownload := false
	if _, err := os.Stat(destination + filename); err == nil {
		// check version match
		if config.TagName == release.TagName {
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

		// update yaml config with new info
		config.TagName = release.TagName
		err = writeConfig(config)
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
	defer response.Body.Close()

	// create file
	output, err := os.Create(destination + filename)
	if err != nil {
		return err
	}
	defer output.Close()

	// write information to file
	_, err = io.Copy(output, response.Body)
	return err
}
