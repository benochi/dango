package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", pingHandler)
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
