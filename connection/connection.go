package connection

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	username   = ""
	password   = ""
	url        = ""
	offsetUrl  = ""
	maxRetries = 5
	delay      = time.Second * 60
)

func createNewRequest(url string) (*http.Response, error) {
	var req *http.Request
	var resp *http.Response
	var err error

	for retries := 0; retries < maxRetries; retries++ {

		req, err = http.NewRequest("GET", url, nil)
		if err == nil {
			req.Header.Set("Accept", "application/json")
			req.SetBasicAuth(username, password)
			client := &http.Client{}
			resp, err = client.Do(req)

			if err == nil && resp.StatusCode == 200 {
				return resp, nil
			}
		}

		time.Sleep(delay)
	}

	log.Println("Retries exhausted, unable to retrieve picklists: " + url + " res:  " + resp.Status)

	if err != nil {
		return nil, err
	}

	if resp != nil {
		return nil, errors.New(resp.Status)
	}

	return nil, errors.New("unknown error occurred")
}

func processRequest(resp *http.Response) ([]byte, error) {

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println("Error while retrieving data from picqer")
	}

	return body, err
}
