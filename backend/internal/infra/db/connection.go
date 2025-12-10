package db

import (
	"context"
	"devhub-backend/internal/config"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	// otelsqlx using v1.10.0
)

type dbContextKey struct{}

func FromContext(ctx context.Context) *sqlx.DB {
	return ctx.Value(dbContextKey{}).(*sqlx.DB)
}

func NewContext(ctx context.Context, db *sqlx.DB) context.Context {
	return context.WithValue(ctx, dbContextKey{}, db)
}

func MustConnect(cfg *config.Config) *sqlx.DB {
	db, err := Connect(cfg)
	if err != nil {
		log.Fatalln("failed to connect to DB:", err)
		return nil
	}
	return db
}

func Connect(cfg *config.Config) (db *sqlx.DB, err error) {
	// Without tracing
	db, err = sqlx.Open("postgres", cfg.DB.URL)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB (sqlx.Open): %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DB.ConnMaxIdleTime)

	return db, nil
}
