package database

import (
	"fmt"

	"github.com/reidaa/ano/internal/database/anime"
	"github.com/reidaa/ano/internal/database/timeseries"
	"github.com/reidaa/ano/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

// Connects to the database using the provided DSN.
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	utils.Info.Println("Connecting to database")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database at %s: %w", dsn, err)
	}

	return db, nil
}

// Prepares the database by migrating the necessary tables.
func Prepare(db *gorm.DB) error {
	utils.Info.Println("Migrating the database")

	err := db.AutoMigrate(&timeseries.TimeseriesModel{}, &anime.AnimeModel{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}
