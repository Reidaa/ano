package jikan

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

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

func (j *Jikan) GetTopAnime(page int, animeType string) (*TopAnimeResponse, error) {
	var responseObj TopAnimeResponse
	params := url.Values{}

	base, err := url.Parse(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %s -> %w", BaseURL, err)
	}

	base.Path += "/top/anime"
	params.Add("page", strconv.Itoa(page))

	if animeType != "" {
		params.Add("type", animeType)
	}

	base.RawQuery = params.Encode()
	url := base.String()

	responseData, err := j.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request %s: %w", url, err)
	}

	err = json.Unmarshal(responseData, &responseObj)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json data: %w", err)
	}

	// To prevent -> 429 Too Many Requests
	time.Sleep(COOLDOWN)

	return &responseObj, nil
}

func (j *Jikan) GetAnimeByID(malID int) (*AnimeResponse, error) {
	var responseObj AnimeResponse

	url := fmt.Sprintf("%s/anime/%d", BaseURL, malID)

	responseData, err := j.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request %s: %w", url, err)
	}

	err = json.Unmarshal(responseData, &responseObj)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json data: %w", err)
	}

	// To prevent -> 429 Too Many Requests
	time.Sleep(COOLDOWN)

	return &responseObj, nil
}
