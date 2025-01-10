package jikan

import (
	"log/slog"
	"net/http"
	"time"
)

type Repository struct {
	client *http.Client
	log    *slog.Logger
}

func NewRepository(logger *slog.Logger) (*Repository, error) {
	r := &Repository{}

	if logger == nil {
		r.log = slog.Default()
	} else {
		r.log = logger
	}

	r.client = &http.Client{
		Timeout: time.Second * Timeout,
	}

	return r, nil
}

func (r *Repository) GetTopAnime() (*TopAnimeResponse, error) {
	return nil, nil
}
