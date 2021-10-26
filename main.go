package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Nosolok/vt/modules"
	"github.com/Nosolok/vt/vtApi"
)

func main() {
	var args = os.Args
	var path string

	if len(args[1:]) != 1 {
		log.Print("Use: vt <path to directory>")
		os.Exit(1)
	} else {
		path = os.Args[1]
	}

	api := modules.GetApiKey()

	var files = modules.Find(path)

	for _, file := range files {
		vtReport := vtApi.FileReport(fmt.Sprintf("%x", file.HashSha1), api)

		switch vtReport.ResponseCode {
		case 0:
			fmt.Println(
				file.Filename, " | ",
				"File is not scanned yet",
			)
		case 1:
			fmt.Println(
				file.Filename, " | ",
				vtReport.Sha1, " | ",
				vtReport.ScanDate, " | ",
				vtReport.Positives,
			)
		default:
			fmt.Println(
				file.Filename, " | ",
				"Unknown response",
			)
		}

		fmt.Println()
	}

}
