package main

import (
	"fmt"
	"os"

	"jishe.wtf/assetdownloader/internal/commands"
)

// Version value is set at build time
var Version string = "dev"

// BuildDate value is set at build time
var BuildDate string = "unknown"

func main() {
	switch os.Args[1] {
	case "list": // list all saved records
		list, err := commands.ListRecords()
		if err != nil {
			fmt.Printf("Error listing records: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("All Saved Records: \n\n" + list)
	case "download":
		err := commands.Download(os.Args[1:])
		if err != nil {
			fmt.Printf("Error downloading asset(s): %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Println("AssetDownloader Version: " + Version + "\nBuild Date: " + BuildDate)
	}
}
