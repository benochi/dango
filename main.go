package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", pingHandler)
	http.ListenAndServe(":5027", nil)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Pinging the server"}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
