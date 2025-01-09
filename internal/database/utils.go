package database

import (
	"fmt"
	"log/slog"

	"github.com/reidaa/ano/internal/database/anime"
	"github.com/reidaa/ano/internal/database/timeseries"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

// Connects to the database using the provided DSN.
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	slog.Info("Connecting to database")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database at %s: %w", dsn, err)
	}

	slog.Info("Connecting to database. Done")
	return db, nil
}

// Prepares the database by migrating the necessary tables.
func Prepare(db *gorm.DB) error {
	slog.Info("Migrating the database")

	err := db.AutoMigrate(&timeseries.TimeseriesModel{}, &anime.AnimeModel{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	slog.Info("Migrating the database. Done")
	return nil
}
