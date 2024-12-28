package database

import (
	"fmt"
	"time"

	"github.com/reidaa/ano/pkg/jikan"
	"github.com/reidaa/ano/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	MAX_BATCH = 100
)

type Database struct {
	client *gorm.DB
}

func New(dbURL string) (*Database, error) {
	n := &Database{}

	db, err := Connect(dbURL)
	if err != nil {
		utils.Error.Println(err)
		return nil, fmt.Errorf("failed to connect to database -> %w", err)
	}

	err = Prepare(db)
	if err != nil {
		utils.Error.Println(err)
		return nil, fmt.Errorf("failed to prepare the database -> %w", err)
	}

	n.client = db

	return n, nil
}

func (db *Database) UpsertTrackedAnimes(animes []jikan.Anime) {
	var data []TrackedModel

	for _, v := range animes {
		data = append(data, TrackedModel{
			MalID:    v.MalID,
			Title:    v.Titles[0].Title,
			ImageURL: v.Images.Jpg.ImageURL,
			Rank:     v.Rank,
			Type:     "anime",
		})
	}

	for _, d := range data {
		utils.Debug.Printf("Upserting in database: %s", d.Title)
		db.client.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "mal_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"title", "image_url", "rank"}),
		})
	}
}

func (db *Database) RetrieveTrackedAnimes() []TrackedModel {
	var result []TrackedModel

	db.client.Find(&result)

	return result
}

// Inserts a list of anime data into the database.
func (db *Database) InsertAnimes(animes []jikan.Anime) {
	var data []animeModel
	var now time.Time = time.Now()

	for _, v := range animes {
		data = append(data, animeModel{
			Timestamp:  now,
			MalID:      v.MalID,
			Title:      v.Titles[0].Title,
			Type:       v.Type,
			Rank:       v.Rank,
			Score:      v.Score,
			ScoredBy:   v.ScoredBy,
			Popularity: v.Popularity,
			Members:    v.Members,
			Favorites:  v.Favorites,
		})
	}

	utils.Info.Println("Writing data to database")
	db.client.CreateInBatches(&data, MAX_BATCH)
	utils.Info.Println("Finished writing data to database")
}
