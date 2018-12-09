package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

const (
	apiServeHost = "localhost"
	apiServePort = "1337"
)

type Status struct {
	CurrentConnectingCount int      `json:"current_connection_count"`
	CurrReqRate            float64  `json:"current_request_rate"`
	ProcessedRequest       int      `json:"processed_request_count"`
	CurrentRequestCount    int      `json:"current_request_count"`
	RemainingJobs          int      `json:"remaining_jobs"`
	CurrentConnecting      []string `json:"current_connected_client"`
	CurrentGoRoutine       int      `json:"current_goroutine_count"`
}

func startStateServer() {
	http.HandleFunc("/status", StatusController)

	fmt.Printf("HTTP Server listening on %s:%s\n", apiServeHost, apiServePort)
	err := http.ListenAndServe(apiServeHost+":"+apiServePort, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func StatusController(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	stat := &Status{
		CurrentConnectingCount: len(connectingClients),
		CurrReqRate:            float64(currentRequestCount) / float64(ratelimitPerSecond),
		ProcessedRequest:       processedRequest,
		CurrentRequestCount:    currentRequestCount,
		RemainingJobs:          len(queryStrings),
		CurrentConnecting:      connectingClients,
		CurrentGoRoutine:       runtime.NumGoroutine(),
	}
	mu.RUnlock()

	b, err := json.Marshal(stat)
	if err != nil {
		fmt.Printf("json marshal failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("json marshal failed"))
		return
	}

	fmt.Fprintf(w, "%s", b)
}
