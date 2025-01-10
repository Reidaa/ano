package anime

import "gorm.io/gorm"

// Struct representing a tracked anime record in the database.
type AnimeModel struct {
	gorm.Model
	Title    string `gorm:"unique"`
	ImageURL string `gorm:"unique;column:image_url"`
	Type     string
	MalID    int `gorm:"unique;column:mal_id"`
	Rank     int
}

func (AnimeModel) TableName() string {
	return "tracked"
}
