package jikan

import (
	"encoding/json"
	"fmt"

	"github.com/reidaa/ano/pkg/utils/netclient"
)

type INet interface {
	Get(url string) ([]byte, error)
}

type Jikan struct {
	http INet
}

func New() (*Jikan, error) {
	var n Jikan = Jikan{}

	n.http = netclient.New()

	return &n, nil
}

func (j *Jikan) GetTopAnime(page int, animeType string, limit int) (*TopAnimeResponse, error) {
	var responseObj TopAnimeResponse
	var err error
	req := GetTopAnimeQuery{
		Page:  page,
		T:     animeType,
		Limit: limit,
	}
	URL, _ := req.URL()

	responseData, err := j.http.Get(URL)
	if err != nil {
		return nil, fmt.Errorf("failed to request %s -> %w", URL, err)
	}

	err = json.Unmarshal(responseData, &responseObj)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json data -> %w", err)
	}

	// if responseObj.Data == nil {
	// 	return nil, fmt.Errorf("failed to unmarshal json data: responseObj.Data is nil")
	// }

	return &responseObj, nil
}

func (j *Jikan) GetAnimeByID(malID int) (*AnimeResponse, error) {
	var responseObj AnimeResponse

	URL := fmt.Sprintf("%s/anime/%d", Endpoint, malID)

	responseData, err := j.http.Get(URL)
	if err != nil {
		return nil, fmt.Errorf("failed to request %s -> %w", URL, err)
	}

	err = json.Unmarshal(responseData, &responseObj)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json data -> %w", err)
	}

	return &responseObj, nil
}
