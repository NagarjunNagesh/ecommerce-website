package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func OKResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode OK response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		log.Printf("failed to encode error response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
