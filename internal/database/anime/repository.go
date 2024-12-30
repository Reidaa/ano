package anime

import (
	"fmt"

	"github.com/reidaa/ano/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// func (r *Repository) Create(anime animeModel) (*animeModel, error)

func (r *Repository) List() ([]*AnimeModel, error) {
	var err error
	animes := make([]*AnimeModel, 0)

	err = r.db.Find(&animes).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all anime -> %w", err)
	}

	return animes, nil
}

func (r *Repository) Read() ([]*AnimeModel, error) {
	var err error
	animes := make([]*AnimeModel, 0)

	err = r.db.Find(&animes).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all anime -> %w", err)
	}

	return animes, nil
}

func (r *Repository) ReadByID(ID uint) (*AnimeModel, error) {
	var err error
	var anime AnimeModel

	err = r.db.Where("id = ?", ID).First(&anime).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve anime with ID: %d -> %w", ID, err)
	}

	return &anime, nil
}

func (r *Repository) ReadByMalID(malID uint) (*AnimeModel, error) {
	var err error
	var anime AnimeModel

	err = r.db.Where("mal_id = ?", malID).First(&anime).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve anime with ID: %d -> %w", malID, err)
	}

	return &anime, nil
}

// func (r *Repository) Update(ID uint) (int, error)
// func (r *Repository) UpdateMany(IDs []uint) (int, error)

func (r *Repository) Upsert(anime *AnimeModel) error {
	utils.Debug.Printf("Upserting in database: %s", anime.Title)

	result := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "mal_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "image_url", "rank", "type"}),
	}).Create(anime)
	if result.Error != nil {
		return fmt.Errorf("failed to upsert anime %s -> %w", anime.Title, result.Error)
	}

	return nil
}

// func (r *Repository) UpsertMany(animes []*AnimeModel) error
// func (r *Repository) Delete(ID uint) (int, error)
