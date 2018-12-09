package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	apiEndPoint = "http://localhost:9090/hello/"
)

var (
	currentRequestCount int
)

type ExternalAPIResponse struct {
	StatusCode int
	Body       string
}

func requestExternalAPI(q string) ExternalAPIResponse {
	timeout := time.Duration(2 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	mu.Lock()
	currentRequestCount++
	mu.Unlock()
	resp, err := client.Get(apiEndPoint + q)

	if err != nil {
		log.Println(err)
		return ExternalAPIResponse{StatusCode: 404, Body: ""}
	}

	mu.Lock()
	currentRequestCount--
	mu.Unlock()

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	return ExternalAPIResponse{StatusCode: resp.StatusCode, Body: string(body)}
}
