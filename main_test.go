package main

import (
	"io"
	"net/http"
	"testing"
)

func TestPingHandler(t *testing.T) {
	resp, err := http.Get("http://localhost:5027/ping")
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
		return
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
		return
	}

	t.Log("Response:", string(body))
}
