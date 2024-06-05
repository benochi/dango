package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

func TestTestHandler(t *testing.T) {
	resp, err := http.Get("http://localhost:5027/test")
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

func TestFileUploadHandler(t *testing.T) {
	// Get absolute path of test file
	testFile := getAbsPath("testfiles/test.csv")

	// Prepare test file
	file, err := os.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Create a new HTTP request with the test file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()

	req := httptest.NewRequest("POST", "http://localhost:5027/file-upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute the request
	rec := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rec, req)

	// Check the response status code
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestFileHandler(t *testing.T) {
	// Get absolute path of test file
	testFile := getAbsPath("testfiles/test.csv")

	// Prepare test file
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Write some content to the test file
	content := []byte("This is a test file")
	_, err = file.Write(content)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP request to retrieve the test file
	req := httptest.NewRequest("GET", "http://localhost:5027/file/test.csv", nil)

	// Execute the request
	rec := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rec, req)

	// Check the response status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Check the response body content
	expectedContent := string(content)
	actualContent := rec.Body.String()
	if actualContent != expectedContent {
		t.Errorf("Expected content %s, got %s", expectedContent, actualContent)
	}
}

// Function to get absolute path of a file
func getAbsPath(relativePath string) string {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		panic(err)
	}
	return absPath
}
