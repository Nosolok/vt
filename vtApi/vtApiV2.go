package vtApi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type fileReportResponse struct {
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
func FileReport(hash string, key string) fileReportResponse {
	var response *http.Response
	var err error

	for done := false; !done; {
		response, err = http.PostForm(
			apiFileReport,
			url.Values{"apikey": {key}, "resource": {hash}},
		)
		defer response.Body.Close()
		if err != nil {
			log.Fatal("error PostForm")
		}

		if response.StatusCode == 200 {
			done = true
		} else {
			time.Sleep(60)
		}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("error ReadAll")
	}

	var report fileReportResponse
	errJson := json.Unmarshal(body, &report)
	if errJson != nil {
		log.Println(errJson)
	}

	return report
}
