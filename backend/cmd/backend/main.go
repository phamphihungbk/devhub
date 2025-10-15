package cmd

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/phamphihungbk/devhub-backend/internal/config"
	"github.com/phamphihungbk/devhub-backend/internal/db"
	"github.com/phamphihungbk/devhub-backend/internal/http/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	err = db.ConnectDB(cfg.DBHost, "5432", cfg.DBUser, cfg.DBPass, "devhub")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	r := gin.Default()
	router.RegisterRoutes(r)

	r.Run(":8080")
}
