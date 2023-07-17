package httpConnection

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type HttpConnection interface {
	CreateNewRequest(url string)
	ProcessRequest(resp *http.Response)
}

type HttpConn struct {
	Response *http.Response
	result   []byte
	Request  *http.Request
	Error    error
}

const (
	username   = ""
	password   = ""
	url        = ""
	offsetUrl  = ""
	maxRetries = 5
	delay      = time.Second * 60
)

func (conn *HttpConn) CreateNewRequest(url string) error {
	//var req *http.Request
	//var resp *http.Response
	//var err error

	for retries := 0; retries < maxRetries; retries++ {

		conn.Request, conn.Error = http.NewRequest("GET", url, nil)
		if conn.Error == nil {
			conn.Request.Header.Set("Accept", "application/json")
			conn.Request.SetBasicAuth(username, password)
			client := &http.Client{}
			conn.Response, conn.Error = client.Do(conn.Request)

			if conn.Error == nil && conn.Response.StatusCode == 200 {
				return nil
			}
		}

		time.Sleep(delay)
	}

	log.Println("Retries exhausted: " + url + " response: " + conn.Response.Status)

	if conn.Error != nil {
		return conn.Error
	}

	if conn.Response != nil {
		return errors.New(conn.Response.Status)
	}

	return errors.New("unknown error occurred")
}

func ProcessRequest(resp *http.Response) ([]byte, error) {

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println("Error while retrieving data from picqer")
	}

	return body, err
}
