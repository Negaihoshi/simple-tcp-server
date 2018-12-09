package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusController(t *testing.T) {
	request, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	record := httptest.NewRecorder()
	handler := http.HandlerFunc(StatusController)
	handler.ServeHTTP(record, request)

	expected := `{"current_connection_count":0,"current_request_rate":0,"processed_request_count":0,"current_request_count":0,"remaining_jobs":0,"current_connected_client":null,"current_goroutine_count":2}`

	if status := record.Code; status != http.StatusOK {
		t.Errorf("wrong status code: except: %v, got %v\n", http.StatusOK, status)
	}

	if record.Body.String() != expected {
		t.Errorf("unexpected body: expect: %v, got %v\n", expected, record.Body.String())
	}
}
