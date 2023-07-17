package http

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type HttpConnection interface {
	CreateNewRequest(url string) error
	Error() error
	Result() []byte
}

type HttpConn struct {
	Response *http.Response
	Request  *http.Request
	result   []byte

	err        error
	Username   string
	Password   string
	Hostname   string
	MaxRetries int
	Delay      time.Duration
}

func (c *HttpConn) CreateNewRequest(url string) error {

	for retries := 0; retries < c.MaxRetries; retries++ {

		c.Request, c.err = http.NewRequest("GET", url, nil)
		if c.err == nil {
			c.Request.Header.Set("Accept", "application/json")
			c.Request.SetBasicAuth(c.Username, c.Password)
			client := &http.Client{}
			c.Response, c.err = client.Do(c.Request)

			if c.err == nil && c.Response.StatusCode == 200 {
				return nil
			}
		}

		time.Sleep(c.Delay)
	}

	log.Println("Retries exhausted: " + url + " response: " + c.Response.Status)

	if c.err != nil {
		return c.err
	}

	if c.Response != nil {
		return errors.New(c.Response.Status)
	}

	c.processRequest()
	if c.processRequest() != nil {
		return c.err
	}

	return errors.New("unknown error occurred")
}

func (c *HttpConn) processRequest() error {
	defer c.Response.Body.Close()

	c.result, c.err = io.ReadAll(c.Response.Body)
	if c.err != nil {
		return c.err
	}

	return nil
}

func (c *HttpConn) Error() error {
	return c.err
}

func (c *HttpConn) Result() []byte {
	return c.result
}
