package httpconn

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type HttpConnection interface {
	CreateNewRequest(url string, methodAndJson ...interface{}) error
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
	Bearer     string
	Hostname   string
	MaxRetries int
	Delay      time.Duration
	Debug      bool
}

func (c *HttpConn) CreateNewRequest(url string, methodAndJson ...interface{}) error {

	//payload := strings.NewReader("{\n  \"trackingNumber\": \"S24DEMO456393\",\n  \"shipmentReference\": \"c6e4fef4-a816-b68f-4024-3b7e4c5a9f81\",\n  \"originCountryCode\": \"CN\",\n  \"destinationCountryCode\": \"US\",\n  \"destinationPostCode\": \"94901\",\n  \"shippingDate\": \"2021-03-01T11:09:00.000Z\",\n  \"courierCode\": [\n    \"us-post\"\n  ],\n  \"courierName\": \"USPS Standard\",\n  \"trackingUrl\": \"https://tools.usps.com/go/TrackConfirmAction?tLabels=S24DEMO456393\",\n  \"orderNumber\": \"DF14R2022\"\n}")

	if c.Debug {
		log.Println("Creating new request: " + url)
	}

	var method string
	var jsonB []byte
	var jsonS *strings.Reader

	if len(methodAndJson) >= 1 {
		if m, ok := methodAndJson[0].(string); ok {
			method = m
		} else {
			return errors.New("invalid method type")
		}
	} else {
		method = "GET" // Default method if not provided
	}

	if len(methodAndJson) >= 2 {
		if j, ok := methodAndJson[1].([]byte); ok {
			jsonB = j
		} else if k, ok := methodAndJson[1].(*strings.Reader); ok {
			jsonS = k
		} else {
			return errors.New("invalid json type")
		}
	}

	for retries := 0; retries < c.MaxRetries; retries++ {

		//c.Request, c.err = http.NewRequest(method, url, bytes.NewBuffer(json))
		//println(bytes.NewBuffer(json))
		if jsonS != nil {
			c.Request, c.err = http.NewRequest(method, url, jsonS)
		} else {
			c.Request, c.err = http.NewRequest(method, url, bytes.NewBuffer(jsonB))
		}

		if c.err == nil {
			c.Request.Header.Set("Accept", "application/json")

			if c.Bearer != "" {
				c.Request.Header.Add("Content-Type", "application/json; charset=utf-8")
				c.Request.Header.Add("Authorization", "Bearer "+c.Bearer)
			}

			if c.Username != "" || c.Password != "" {
				c.Request.SetBasicAuth(c.Username, c.Password)
			}

			client := &http.Client{}
			c.Response, c.err = client.Do(c.Request)

			if c.err == nil && (c.Response.StatusCode == 200 || c.Response.StatusCode == 201) {

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
