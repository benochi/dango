package main

import (
	"fmt"
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
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
		return
	}

	fmt.Println("Response:", string(body))
}
