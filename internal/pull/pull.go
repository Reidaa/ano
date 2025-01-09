package pull

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/reidaa/ano/internal/database"
	"github.com/reidaa/ano/internal/database/anime"
	"github.com/reidaa/ano/internal/database/timeseries"
	"github.com/reidaa/ano/pkg/jikan"
	"github.com/reidaa/ano/pkg/utils/intset"
)

type Config struct {
	DatabaseURL   string
	Top           int
	SkipRetrieval bool
}

type Scrapper struct {
	animeRepository *anime.Repository
	tsRepository    *timeseries.Repository
	log             *slog.Logger
	animes          intset.IntSet
	conf            Config
	useDatabase     bool
}

func New(conf Config, logger *slog.Logger) (*Scrapper, error) {
	s := &Scrapper{
		conf:        conf,
		animes:      *intset.New(),
		useDatabase: false,
		log:         logger,
	}

	if s.conf.DatabaseURL != "" {
		s.useDatabase = true
	}

	return s, nil
}

func (s *Scrapper) Start() error {
	var tops []jikan.Anime
	var err error

	if s.useDatabase {
		err = s.connectToDatabase()
		if err != nil {
			return err
		}
	}

	tops, err = s.checkTop()
	if err != nil {
		return err
	}

	if s.useDatabase {
		for i := range tops {
			err = s.animeRepository.Upsert(&anime.AnimeModel{
				MalID:    tops[i].MalID,
				Title:    tops[i].Titles[0].Title,
				ImageURL: tops[i].Images.Jpg.ImageURL,
				Rank:     tops[i].Rank,
				Type:     tops[i].Type,
			})
			if err != nil {
				return fmt.Errorf("failed to upsert anime -> %w", err)
			}
		}
		if !s.conf.SkipRetrieval {
			err = s.retrieveAnimeFromDB()
			if err != nil {
				return err
			}
		}
	}

	data := s.getAnimeData(s.animes.Slice())

	if s.useDatabase {
		now := time.Now()
		for i := range data {
			d := &timeseries.TimeseriesModel{
				Timestamp:  now,
				MalID:      data[i].MalID,
				Title:      data[i].Titles[0].Title,
				Type:       data[i].Type,
				Rank:       data[i].Rank,
				Score:      data[i].Score,
				ScoredBy:   data[i].ScoredBy,
				Popularity: data[i].Popularity,
				Members:    data[i].Members,
				Favorites:  data[i].Favorites,
			}
			err = s.tsRepository.Create(d)
			if err != nil {
				return fmt.Errorf("failed to insert timeseries data -> %w", err)
			}
		}
	}

	s.render(data)

	return nil
}

func (s *Scrapper) checkTop() ([]jikan.Anime, error) {
	slog.Info(fmt.Sprintf("Checking the top %d anime", s.conf.Top))

	tops, err := jikan.TopAnimeByRank(s.conf.Top)
	if err != nil {
		return nil, fmt.Errorf("failed retrieve the top %d anime -> %w", s.conf.Top, err)
	}

	for i := range tops {
		s.animes.Insert(tops[i].MalID)
	}

	slog.Info(fmt.Sprintf("Checking the top %d anime. Done", s.conf.Top))
	return tops, nil
}

func (s *Scrapper) connectToDatabase() error {
	slog.Info("Establishing connection to database")
	db, err := database.Connect(s.conf.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database -> %w", err)
	}

	err = database.Prepare(db)
	if err != nil {
		return fmt.Errorf("failed to prepare the database -> %w", err)
	}

	s.animeRepository = anime.New(db)
	s.tsRepository = timeseries.New(db)

	slog.Info("Establishing connection to database. Done")
	return nil
}

func (s *Scrapper) retrieveAnimeFromDB() error {
	tracked, err := s.animeRepository.Read()
	if err != nil {
		return fmt.Errorf("failed to retrieve anime from database -> %w", err)
	}

	for i := range tracked {
		s.animes.Insert(tracked[i].MalID)
	}

	return nil
}

func (s *Scrapper) getAnimeData(malIDs []int) []jikan.Anime {
	slog.Info(fmt.Sprintf("Gathering data about the %d entries", len(malIDs)))

	var data []jikan.Anime

	for i := range malIDs {
		d, err := jikan.AnimeByID(malIDs[i])
		jikan.Cooldown()
		if err != nil {
			slog.Warn("failed to fetch data, skipping", slog.Any("error", err))
		} else {
			data = append(data, *d)
		}
	}

	slog.Info(fmt.Sprintf("Gathering data. Done"))
	return data
}

func (s *Scrapper) render(animes []jikan.Anime) {
	t := table.NewWriter()

	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"mal_id", "title", "rank", "score", "members", "favorites"})

	for i := range animes {
		t.AppendRow(table.Row{
			animes[i].MalID, animes[i].Titles[0].Title, animes[i].Rank, animes[i].Score, animes[i].Members, animes[i].Favorites,
		})
	}

	t.Render()
}
