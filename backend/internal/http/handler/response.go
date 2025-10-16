package handler

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ErrorResponse writes an error message as JSON
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, map[string]string{"error": message})
}
