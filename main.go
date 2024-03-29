package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Nosolok/vt/modules"
	"github.com/Nosolok/vt/vtApi"

	_ "github.com/mattn/go-sqlite3"
)

const maxFileSize = 1024 * 1024 * 32

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

	dbFile := "./vt.db"

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = os.Stat(dbFile)
	if os.IsNotExist(err) {
		modules.InitDatabase(db)
	}

	var files = modules.Find(path)

	fmt.Println(
		"source", " | ",
		"flename", " | ",
		"sha256", " | ",
		"scan date", " | ",
		"positives", " | ",
		"total",
	)

	for _, file := range files {

		fileInfoDb := modules.CheckHash(db, fmt.Sprintf("%x", file.HashSha256))

		if fileInfoDb.Id == 0 {
			// hash not found in local database, retrieve it from virustotal

			vtReport := vtApi.FileReport(fmt.Sprintf("%x", file.HashSha1), api)

			switch vtReport.ResponseCode {
			case -2:
				// If the requested item is still queued for analysis
				fmt.Println(
					file.Filename, " | ",
					"File in in queue for analysis",
				)
			case 0:
				// the item you searched for was not present in VirusTotal's dataset
				fmt.Println(
					file.Filename, " | ",
					"File is not scanned yet",
				)

				f, err := os.Open(path + fmt.Sprintf("%c", os.PathSeparator) + file.Filename)
				if err != nil {
					log.Fatal("error open file", err)
				}
				defer f.Close()

				var fullFilename = path + fmt.Sprintf("%c", os.PathSeparator) + file.Filename

				fileInfo, err := f.Stat()
				// fmt.Println(fileInfo.Size())

				if fileInfo.Size() > maxFileSize {
					fmt.Println("big file")
					// vtApi.UploadBigFile(f, fullFilename, api)
				} else {
					fmt.Println("small file")
					vtApi.UploadFile(vtApi.Api3FileUpload, f, fullFilename, api)
				}

			case 1:
				// item was indeed present and it could be retrieved
				modules.StoreCheck(db, vtReport)
				fmt.Println(
					"vt", " | ",
					file.Filename, " | ",
					vtReport.Sha256, " | ",
					vtReport.ScanDate, " | ",
					vtReport.Positives, " | ",
					vtReport.Total,
				)
			default:
				fmt.Println(
					file.Filename, " | ",
					"Unknown response",
				)
			}
			fmt.Println()
		} else {
			// hash found in local database, output cached information

			fmt.Println(
				"db", " | ",
				file.Filename, " | ",
				fileInfoDb.Sha256, " | ",
				fileInfoDb.Scan_date, " | ",
				fileInfoDb.Positives, " | ",
				fileInfoDb.Total,
			)
		}

	}

}
