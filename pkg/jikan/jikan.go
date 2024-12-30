package jikan

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

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
	query := url.Values{}

	if page > 0 {
		query.Add("page", strconv.Itoa(page))
	}

	if limit > 0 {
		query.Add("limit", strconv.Itoa(limit))
	}

	if animeType != "" {
		query.Add("type", animeType)
	}

	base, err := url.Parse(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %s -> %w", BaseURL, err)
	}

	base.Path += "/top/anime"
	base.RawQuery = query.Encode()
	URL := base.String()

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

	URL := fmt.Sprintf("%s/anime/%d", BaseURL, malID)

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
