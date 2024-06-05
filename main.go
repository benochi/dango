package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/file-upload", fileUploadHandler)
	http.HandleFunc("/file/", fileHandler)

	fmt.Println("Server running on PORT:5027")
	err := http.ListenAndServe(":5027", nil)
	if err != nil {
		panic(err)
	}

}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Pinging the server"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "This is the test handler route."}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Ensure the "uploads" directory exists
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to create uploads directory", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file on the server
	f, err := os.Create("./uploads/" + handler.Filename)
	if err != nil {
		http.Error(w, "Failed to create file on server", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the uploaded file to the newly created file
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "Failed to save file on server", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message":  "File was successfully uploaded!",
		"filename": handler.Filename,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[len("/file/"):]
	filepath := "./uploads/" + filename

	// Open the requested file
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Copy the file to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
		return
	}
}
