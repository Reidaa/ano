package jikan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"malstat/scrapper/pkg/utils"
)

type topAnimeResponse struct {
	Data       []Anime    `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func topAnime(page int, animeType string) (*topAnimeResponse, error) {
	var responseObj topAnimeResponse
	url := fmt.Sprintf("%s/top/anime?page=%d", Base_url, page)

	if animeType != "" {
		url = fmt.Sprintf("%s&type=%s", url, animeType)
	}

	utils.Debug.Printf("GET %s", url)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, errors.New(response.Status)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseData, &responseObj)
	if err != nil {
		return nil, err
	}

	// To prevent -> 429 Too Many Requests
	time.Sleep(Cooldown)

	return &responseObj, nil
}

func TopAnime(n int) (*[]Anime, error) {
	var data []Anime

	types := []string{"tv", "movie", "ova", "tv_special", "special"}

	for t := 0; t != len(types); t++ {
		response, err := topAnime(1, types[t])
		if err != nil {
			return nil, err
		}

		data = append(data, response.Data...)

		for i := 2; i <= n/response.Pagination.Items.PerPage; i++ {
			response, err := topAnime(i, types[t])
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

	for i := 1; maxCurrentRank < maxRank; i++ {
		response, err := topAnime(i, "")
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
