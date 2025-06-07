package main

import (
	"encoding/json"
	"io"

	"fmt"
	"net/http"
	"os"

	"jishe.wtf/assetdownloader/internal/records"
)

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

var destination string
var finalArg string

var hasDownloadAllFlag bool

func main() {
	usageString := "Usage: <repo> <filename> [<repo> <filename>] <destination> <-a | --download-all>\n"
	usageString += "<repo> <filename> <destination> [<repo> <filename> <destination>]"

	// check for download all flag
	finalArg = os.Args[len(os.Args)-1]
	if finalArg == "-a" || finalArg == "--download-all" {
		hasDownloadAllFlag = true
	}

	// check if minimum number of arguments is provided and
	// if amount of arguments is NOT a value of the below sequence
	// 4, 7, 10, 13, 16, 19, 22, 25, ...
	if ((len(os.Args) - 1) < 3) || ((len(os.Args)-1)%3 != 0 && !hasDownloadAllFlag) {
		fmt.Println("Error: Invalid number of arguments provided.")
		fmt.Println(usageString)
		os.Exit(1)
	}

	// only run if flag is at the end
	if hasDownloadAllFlag {
		if (len(os.Args)-1)%2 != 0 {
			fmt.Println("Error: Invalid number of arguments provided with download all flag.")
			fmt.Println(usageString)
			os.Exit(1)
		}

		fmt.Println("Flag identified, downloading all files to the same destination.")
		destination = parseDownloadPath(os.Args[len(os.Args)-2])
	}

	// loop through all given repo + filename pairs
	for i := 1; i < len(os.Args)-2; i += 2 {
		repo := os.Args[i]
		filename := os.Args[i+1]

		// get download destination for current asset
		if !hasDownloadAllFlag {
			destination = parseDownloadPath(os.Args[i+2])
			i++
		}

		// get latest release, in separate function to isolate defer calls
		release, err := getRelease(repo)
		if err != nil {
			fmt.Printf("Error getting release: %v\n", err)
			os.Exit(1)
		}

		// check if already exists
		record, err := records.Load(repo, filename, destination)
		if err != nil {
			fmt.Printf("Error loading record: %v\n", err)
			os.Exit(1)
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

		// download file
		err = downloadFile(destination, filename, file.DownloadURL)
		if err != nil {
			fmt.Printf("Error downloading file: %v\n", err)
			os.Exit(1)
		}

		// update yaml record with new info
		record.TagName = release.TagName
		record.AuthorName = release.Author.Name

		err = records.Write(record)
		if err != nil {
			fmt.Printf("Error writing config file: %v\n", err)
			os.Exit(1)
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

// setting destination variable and path formatting
func parseDownloadPath(path string) string {
	// ensure trailing "/" always present
	if path[len(path)-1:] != "/" {
		path += "/"
	}

	return path
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
