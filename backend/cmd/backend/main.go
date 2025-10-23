package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/phamphihungbk/devhub-backend/internal/config"
	"github.com/phamphihungbk/devhub-backend/internal/http/router"
)

func main() {
	_, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// db := db.Connect(cfg)

	r := gin.Default()
	router.RegisterRoutes(r)

	r.Run(":8080")
}
