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
		"flename", " | ",
		"sha256", " | ",
		"scan date", " | ",
		"positives", " | ",
		"total",
	)

	for _, file := range files {

		fileInfoDb := modules.CheckHash2(db, fmt.Sprintf("%x", file.HashSha256))

		if fileInfoDb.Id == 0 {
			// hash not found in local database, retrieve it from virustotal

			vtReport := vtApi.FileReport(fmt.Sprintf("%x", file.HashSha1), api)

			switch vtReport.ResponseCode {
			case 0:
				// the item you searched for was not present in VirusTotal's dataset
				fmt.Println(
					file.Filename, " | ",
					"File is not scanned yet",
				)
			case 1:
				// item was indeed present and it could be retrieved
				modules.StoreCheck(db, vtReport)
				fmt.Println(
					file.Filename, " | ",
					vtReport.Sha256, " | ",
					vtReport.ScanDate, " | ",
					vtReport.Positives, " | ",
					fileInfoDb.Total,
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
				file.Filename, " | ",
				fileInfoDb.Sha256, " | ",
				fileInfoDb.Scan_date, " | ",
				fileInfoDb.Positives, " | ",
				fileInfoDb.Total,
			)
		}

	}

}
