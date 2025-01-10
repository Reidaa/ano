package jikan

import (
	"fmt"
	"sort"
	"time"
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

func TopAnimeByRank(maxRank int) ([]Anime, error) {
	var data []Anime
	var maxCurrentRank int = 0
	var limit int = DefaultLimit

	if maxRank < DefaultLimit {
		limit = maxRank
	}

	client, err := New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize jikan client -> %w", err)
	}

	for i := 1; maxCurrentRank < maxRank; i++ {
		response, err := client.GetTopAnime(i, "", limit)
		if err != nil {
			return nil, err
		}

		// To prevent -> 429 Too Many Requests
		time.Sleep(CooldownDuration)

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

// To prevent -> 429 Too Many Requests.
func Cooldown() {
	time.Sleep(CooldownDuration)
}
