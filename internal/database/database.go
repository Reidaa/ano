package database

import (
	"fmt"
	"time"

	"github.com/reidaa/ano/internal/database/anime"
	"github.com/reidaa/ano/internal/database/timeseries"
	"github.com/reidaa/ano/pkg/jikan"
	"github.com/reidaa/ano/pkg/utils"
	"gorm.io/gorm"
)

type Database struct {
	client              *gorm.DB
	animeRepository     *anime.Repository
	timeserieRepository *timeseries.Repository
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
	n.animeRepository = anime.New(db)
	n.timeserieRepository = timeseries.New(db)

	return n, nil
}

func (db *Database) UpsertTrackedAnimes(animes []jikan.Anime) error {
	var err error

	for i := range animes {
		err = db.animeRepository.Upsert(&anime.AnimeModel{
			MalID:    animes[i].MalID,
			Title:    animes[i].Titles[0].Title,
			ImageURL: animes[i].Images.Jpg.ImageURL,
			Rank:     animes[i].Rank,
			Type:     animes[i].Type,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert anime -> %w", err)
		}
	}

	return nil
}

func (db *Database) RetrieveTrackedAnimes() ([]*anime.AnimeModel, error) {
	animes, err := db.animeRepository.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read anime from database -> %w", err)
	}

	return animes, nil
}

func (db *Database) InsertAnimes(animes []jikan.Anime) error {
	var err error
	now := time.Now()

	for i := range animes {
		d := &timeseries.TimeseriesModel{
			Timestamp:  now,
			MalID:      animes[i].MalID,
			Title:      animes[i].Titles[0].Title,
			Type:       animes[i].Type,
			Rank:       animes[i].Rank,
			Score:      animes[i].Score,
			ScoredBy:   animes[i].ScoredBy,
			Popularity: animes[i].Popularity,
			Members:    animes[i].Members,
			Favorites:  animes[i].Favorites,
		}
		err = db.timeserieRepository.Create(d)
		if err != nil {
			return fmt.Errorf("failed to insert timeseries data -> %w", err)
		}
	}

	return nil
}
