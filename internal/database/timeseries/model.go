package timeseries

import (
	"time"

	"gorm.io/gorm"
)

// Struct representing an anime record in the database.
type TimeseriesModel struct {
	Timestamp time.Time
	gorm.Model
	Title      string
	Type       string
	MalID      int
	Rank       int
	ScoredBy   int
	Popularity int
	Members    int
	Favorites  int
	Score      float32
}

func (TimeseriesModel) TableName() string {
	return "animes"
}
