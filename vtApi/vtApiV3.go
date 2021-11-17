package vtApi

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

const api3Url string = "https://www.virustotal.com/api/"
const api3Version string = "v3/"
const api3FileUpload string = api3Url + api3Version + "files"

func UploadFile(file io.Reader, filename string, key string) {
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

	req, _ := http.NewRequest("POST", api3FileUpload, &b)

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
