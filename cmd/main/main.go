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
	uString := "Usage: assetdownloader <command> [options]\n\nAvailable commands:\n  list\n  download\n  version"

	if len(os.Args) < 2 {
		fmt.Println(uString)
		os.Exit(1)
	}

	switch os.Args[1] {
	// list all saved records
	case "list":
		list, err := commands.ListRecords()
		if err != nil {
			fmt.Printf("Error listing records: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("All Saved Records: \n\n" + list)
	// download assets from Github repositories
	case "download":
		err := commands.Download(os.Args[1:])
		if err != nil {
			fmt.Printf("Error downloading asset(s): %v\n", err)
			os.Exit(1)
		}
	// display version information
	case "version":
		fmt.Println("AssetDownloader Version: " + Version + "\nBuild Date: " + BuildDate)
	// not a real command
	default:
		fmt.Println("Unknown command: " + os.Args[1] + "\n\n" + uString)
	}
}
