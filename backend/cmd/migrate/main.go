package main

import (
	"flag"
	"log"
	"os"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	// import your config and db packages as needed
)

func main() {
	direction := flag.String("direction", "up", "Migration direction: up or down")
	steps := flag.Int("steps", 1, "Number of steps to migrate (for down)")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db := db.Connect(cfg)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Error creating postgres driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://backend/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Error creating migration instance: %v", err)
	}

	switch *direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migrations applied successfully.")
	case "down":
		if err := m.Steps(-*steps); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Printf("Rolled back %d migration(s) successfully.", *steps)
	default:
		log.Fatalf("Unknown direction: %s", *direction)
	}
}
