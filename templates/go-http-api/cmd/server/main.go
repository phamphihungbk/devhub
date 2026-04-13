package main

import (
	"log"
	"net/http"

	"{{MODULE_PATH}}/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.Health)
	log.Fatal(http.ListenAndServe(":{{PORT}}", mux))
}
