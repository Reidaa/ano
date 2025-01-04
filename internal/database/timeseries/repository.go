package timeseries

import (
	"fmt"

	"github.com/reidaa/ano/pkg/utils"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(data *TimeseriesModel) error {
	utils.Debug.Printf("Inserting in db snapshot of %s taken at %s", data.Title, data.Timestamp.UTC())
	result := r.db.Create(&data)
	if result.Error != nil {
		return fmt.Errorf("failed to insert data -> %w", result.Error)
	}

	return nil
}

// func (r *Repository) List() ([]*TimeseriesModel, error)
// func (r *Repository) Read() ([]*TimeseriesModel, error)
// func (r *Repository) ReadByID(ID uint) (*TimeseriesModel, error)

func (r *Repository) ReadByMalID(malID uint) (*TimeseriesModel, error) {
	var err error
	var data TimeseriesModel

	err = r.db.Where("mal_id = ?", malID).First(&data).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data with ID: %d -> %w", malID, err)
	}

	return &data, nil
}

// func (r *Repository) Update(ID uint) (int, error)
// func (r *Repository) UpdateMany(IDs []uint) (int, error)

// func (r *Repository) Upsert(data *TimeseriesModel) error

// func (r *Repository) UpsertMany(animes []*AnimeModel) error
// func (r *Repository) Delete(ID uint) (int, error)
