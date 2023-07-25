package httpconn

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
	Debug      bool
}

func (c *HttpConn) CreateNewRequest(url string) error {

	if c.Debug {
		log.Println("Creating new request: " + url)
	}

	for retries := 0; retries < c.MaxRetries; retries++ {

		c.Request, c.err = http.NewRequest("GET", url, nil)
		if c.err == nil {
			c.Request.Header.Set("Accept", "application/json")
			c.Request.SetBasicAuth(c.Username, c.Password)
			client := &http.Client{}
			c.Response, c.err = client.Do(c.Request)

			if c.err == nil && c.Response.StatusCode == 200 {

				if c.Debug {
					log.Println("Response 200")
				}

				if c.processRequest() != nil {
					return c.err
				}

				if c.Debug {
					log.Println("Request succesfully executed")
				}

				return nil
			}
		}

		if c.Debug {
			log.Println("Attempt failed, sleeping for ", c.Delay)
		}

		time.Sleep(c.Delay)
	}

	if c.Debug {
		log.Println("Retries exhausted")
	}

	if c.err != nil {
		return c.err
	}

	if c.Response != nil {
		return errors.New(c.Response.Status)
	}

	return errors.New("unknown error occurred")
}

func (c *HttpConn) processRequest() error {
	defer c.Response.Body.Close()

	if c.Debug {
		log.Println("Start Reading body")
	}

	c.result, c.err = io.ReadAll(c.Response.Body)
	if c.err != nil {
		return c.err
	}

	if c.Debug {
		log.Println("Reading body succesfull")
	}

	return nil
}

func (c *HttpConn) Error() error {
	return c.err
}

func (c *HttpConn) Result() []byte {
	//fmt.Printf("c.result: %v\n", string(c.result))
	return c.result
}
