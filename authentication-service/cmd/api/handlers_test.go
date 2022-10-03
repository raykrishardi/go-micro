package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test authenticate by simulating a POST request from the browser
func TestAuthenticate(t *testing.T) {
	postBody := map[string]interface{}{
		"email":    "me@here.com",
		"password": "verysecret",
	}

	body, _ := json.Marshal(postBody)

	req, _ := http.NewRequest("POST", "/authenticate", bytes.NewReader(body))
	// Response recorder
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.Authenticate)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected http.StatusAccepted but go %d", rr.Code)
	}
}
