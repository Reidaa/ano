package jikan

import (
	"fmt"
	"sort"
	"time"

	"github.com/reidaa/ano/pkg/utils"
)

const (
	MaxAllowedHitPerDay int           = 60 * 60 * 24
	MaxSafeHitPerDay    int           = 60 * 60 * 20
	BaseURL             string        = "https://api.jikan.moe/v4"
	COOLDOWN            time.Duration = time.Second
)

func RemoveUnrankedAnime(in []Anime) []Anime {
	var out []Anime

	for i := 0; i != len(in); i++ {
		if in[i].Rank != 0 {
			out = append(out, in[i])
		}
	}

	return out
}

func TopAnime(n int) (*[]Anime, error) {
	var data []Anime

	client, err := New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize jikan client -> %w", err)
	}

	types := []string{"tv", "movie", "ova", "tv_special", "special"}

	for t := 0; t != len(types); t++ {
		response, err := client.GetTopAnime(1, types[t])
		if err != nil {
			return nil, err
		}

		data = append(data, response.Data...)

		for i := 2; i <= n/response.Pagination.Items.PerPage; i++ {
			response, err := client.GetTopAnime(i, types[t])
			if err != nil {
				return nil, err
			}
			data = append(data, response.Data...)
		}
	}

	for i := 0; i != len(data); i++ {
		utils.Debug.Println(data[i].Titles[0].Title, data[i].Rank)
	}

	return &data, nil
}

func TopAnimeByRank(maxRank int) ([]Anime, error) {
	var data []Anime
	var maxCurrentRank int = 0

	client, err := New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize jikan client -> %w", err)
	}

	for i := 1; maxCurrentRank < maxRank; i++ {
		response, err := client.GetTopAnime(i, "")
		if err != nil {
			return nil, err
		}
		data = append(data, response.Data...)
		maxCurrentRank = response.Data[len(response.Data)-1].Rank
	}

	data = RemoveUnrankedAnime(data)

	sort.Slice(data, func(i, j int) bool {
		return data[i].Score > data[j].Score
	})

	return data, nil
}

func AnimeByID(id int) (*Anime, error) {
	client, err := New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize jikan client -> %w", err)
	}

	d, err := client.GetAnimeByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get anime: id %d -> %w", id, err)
	}

	return &d.Data, nil
}
