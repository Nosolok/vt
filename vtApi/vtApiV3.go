package vtApi

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

type uploadUrl struct {
	Data string `json:"data`
}

const api3Url string = "https://www.virustotal.com/api/"
const api3Version string = "v3/"
const Api3FileUpload string = api3Url + api3Version + "files"
const api3BigFileUpload string = api3Url + api3Version + "files/upload_url"

func UploadFile(url string, file io.Reader, filename string, key string) {
	b := bytes.Buffer{}
	writer := multipart.NewWriter(&b)
	form, err := writer.CreateFormFile("file", filename)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(form, file)
	if err != nil {
		log.Fatal(err)
	}

	writer.Close()

	req, _ := http.NewRequest("POST", url, &b)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Apikey", key)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
}

func UploadBigFile(file io.Reader, filename string, key string) {
	req, _ := http.NewRequest("GET", api3BigFileUpload, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Apikey", key)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data uploadUrl
	errJson := json.Unmarshal(body, &data)
	if errJson != nil {
		log.Println(errJson)
	}

	UploadFile(data.Data, file, filename, key)
}
