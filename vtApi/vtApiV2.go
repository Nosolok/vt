package vtApi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type FileReportResponse struct {
	ResponseCode int    `json:"response_code"`
	VerboseMsg   string `json:"verbose_msg"`
	Resource     string `json:"resource"`
	ScanID       string `json:"scan_id"`
	Md5          string `json:"md5"`
	Sha1         string `json:"sha1"`
	Sha256       string `json:"sha256"`
	ScanDate     string `json:"scan_date"`
	Permalink    string `json:"permalink"`
	Positives    int    `json:"positives"`
	Total        int    `json:"total"`
	Scans        map[string]struct {
		Detected bool   `json:"detected"`
		Version  string `json:"version"`
		Result   string `json:"result"`
		Update   string `json:"update"`
	} `json:"scans"`
}

const apiUrl string = "https://www.virustotal.com/vtapi/"
const apiVersion string = "v2/"
const apiFileReport string = apiUrl + apiVersion + "file/report"

// FileReport get report about already scanned files
// Return report about file
func FileReport(hash string, key string) FileReportResponse {
	var response *http.Response
	var err error
	var isApiLimitExceeded bool = false

	for done := false; !done; {
		response, err = http.PostForm(
			apiFileReport,
			url.Values{"apikey": {key}, "resource": {hash}},
		)
		if err != nil {
			log.Fatal("error PostForm", err)
		}
		defer response.Body.Close()

		switch response.StatusCode {
		case 200:
			// normal response
			isApiLimitExceeded = false
			done = true
		case 204:
			// reponse when API quota limit exceeded. Maybe limitation by minutes
			// so if isApiLimitExceeded appear twice it means daily limit is used up
			if isApiLimitExceeded {
				log.Fatal("API limit exceeded")
			}

			isApiLimitExceeded = true
			time.Sleep(60 * time.Second)
		default:
			log.Println("default switch")
		}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("error ReadAll")
	}

	var report FileReportResponse
	errJson := json.Unmarshal(body, &report)
	if errJson != nil {
		log.Println(errJson)
	}

	return report

}
