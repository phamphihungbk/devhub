package handler

import (
	"encoding/json"
	"net/http"

	"[[ MODULE_PATH ]]/internal/service"
)

func Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(service.Health())
}
