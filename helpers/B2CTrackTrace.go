package helpers

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"
)

const (
	trackingStr       string = "Your tracking number: "
	offsetTrackingStr int    = len(trackingStr)
	notFoundMsg              = "Not found"
)

/*
Scrape the track and trace code from this website https://www.trackyourparcel.eu/
Return empty string if track and trace could not be found and err if website could not be read
*/
func GetTrackAndTrace(website string) (string, error) {
	websiteBody, err := readWebsiteBody(website)
	if err != nil {
		return "", err
	}

	startIndex := strings.Index(websiteBody, trackingStr)
	if startIndex == -1 {
		return "", nil
	}

	var trackAndTraceCode strings.Builder
	for _, r := range websiteBody[startIndex+offsetTrackingStr:] {
		if isValidCharForTrackAndTrace(string(r)) {
			trackAndTraceCode.WriteByte(byte(r))
		} else {
			break
		}
	}

	return trackAndTraceCode.String(), nil
}

func GetCourier(website string) (string, error) {
	websiteBody, err := readWebsiteBody(website)
	if err != nil {
		return "", err
	}

	cntDHL := strings.Count(websiteBody, "DHL")
	cntBpost := strings.Count(websiteBody, "bpost.be")
	cntDPD := strings.Count(websiteBody, "DPD")

	if cntDHL > 0 {
		return "dhl", nil
	}

	if cntBpost > 0 {
		return "be-post", nil
	}

	if cntDPD > 0 {
		return "dpd", nil
	}

	return "", nil
}

func isValidCharForTrackAndTrace(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return true
		}
	}

	return false
}

func readWebsiteBody(website string) (string, error) {

	const maxRetries = 5
	const delay = time.Second * 60

	var resp *http.Response

	for retries := 0; retries < maxRetries; retries++ {

		req, err := http.NewRequest("GET", website, nil)
		if err == nil {

			client := &http.Client{}
			resp, err := client.Do(req)

			if err == nil && resp.StatusCode == 200 {

				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					return string(body), nil
				}
			}
		}

		time.Sleep(delay)
	}

	log.Println("Retries exhausted, unable to read website: ", website)

	return "", errors.New(resp.Status)
}
