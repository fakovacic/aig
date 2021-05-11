package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fakovacic/aig/internal/log"
)

// go run main.go -case= -showResponse=
func main() {
	var (
		caseFile     string
		showResponse string
	)
	flag.StringVar(&caseFile, "case", "", "case filename")
	flag.StringVar(&showResponse, "showResponse", "false", "show response")
	flag.Parse()
	if caseFile == "" {
		log.Error("case file can't be empty")
	}

	caseList, err := NewCase(caseFile)
	if err != nil {
		log.Errorf("new case:%v", err)
	}

	if caseList.URL != "" {
		log.Infof("%s\n", caseList.URL)
	}

	for i := range caseList.Requests {
		if showResponse == "true" {
			caseList.Requests[i].ShowResponse = true
		}

		if caseList.Requests[i].URL != "" {
			log.Infof("%s\n", caseList.Requests[i].URL)
		}

		if caseList.URL != "" {
			caseList.Requests[i].URL = caseList.URL
		}

		if len(caseList.Headers) != 0 {
			if len(caseList.Requests[i].Headers) == 0 {
				caseList.Requests[i].Headers = make(map[string]string)
			}

			for aKey, aVal := range caseList.Headers {
				caseList.Requests[i].Headers[aKey] = aVal
			}
		}

		if caseList.Requests[i].Forever {
			for {
				caseList.Requests[i].Send()
			}
		}

		caseList.Requests[i].Send()
	}
}

// NewCase return new http requests case from input file
func NewCase(file string) (Case, error) {
	var cases Case

	jsonFile, err := os.Open(
		fmt.Sprintf("./cases/%s", file),
	)
	if err != nil {
		return cases, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &cases)
	if err != nil {
		return cases, err
	}
	return cases, nil
}

// Case struct
type Case struct {
	URL      string            `json:"url,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
	Requests []Requests        `json:"requests"`
}

// Requests struct
type Requests struct {
	URL          string            `json:"url,omitempty"`
	Action       string            `json:"action"`
	Method       string            `json:"method"`
	Headers      map[string]string `json:"headers,omitempty"`
	Response     Response          `json:"response,omitempty"`
	Forever      bool              `json:"forever,omitempty"`
	Body         interface{}       `json:"body,omitempty"`
	ShowResponse bool
}

// Send send http requests
func (r *Requests) Send() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*5,
	)
	defer cancel()

	b, _ := json.Marshal(r.Body)

	req, _ := http.NewRequestWithContext(
		ctx,
		r.Method,
		fmt.Sprintf("%s%s", r.URL, r.Action),
		strings.NewReader(string(b)),
	)

	if len(r.Headers) != 0 {
		for key, value := range r.Headers {
			req.Header.Add(key, value)
		}
	}

	showResponse := false
	if r.ShowResponse {
		showResponse = r.ShowResponse
	}

	start := time.Now()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("error:%v", err)
	}

	log.Infof("%s - %s - %d - %s\n",
		r.Action,
		r.Method,
		res.StatusCode,
		time.Since(start).String(),
	)

	if r.Response.Status != 0 && r.Response.Status != res.StatusCode {
		log.Error("- error status code")
	}

	defer res.Body.Close()
	if showResponse {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Errorf("error reading body:%v", err)
		}

		log.Debugf("- response: %s", string(body))
	}
}

// Response struct
type Response struct {
	Status int    `json:"status,omitempty"`
	Body   string `json:"body,omitempty"`
}
