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
	usage := "Usage: assetdownloader <command> [options]\n\nAvailable commands:\n  list\n  download\n  version"

	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

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
	default:
		fmt.Println("Unknown command: " + os.Args[1] + "\n\n" + usage)
	}
}
