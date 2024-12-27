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

// GetTopAnime retrieves a list of top anime from the Jikan API.
//
// Parameters:
//   - page: The page number for pagination (if > 0)
//   - animeType: Filter results by anime type (e.g., "tv", "movie", etc.)
//   - limit: Maximum number of results per page (if > 0)
//
// Returns:
//   - *TopAnimeResponse: Contains the list of top anime and pagination information
//   - error: Non-nil if an error occurred during the request or data processing
//
// The function makes a GET request to the Jikan API's /top/anime endpoint.
// It supports pagination and filtering by anime type. If page or limit are <= 0,
// those query parameters will be omitted from the request.
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
	url := base.String()

	responseData, err := j.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request %s -> %w", url, err)
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

	url := fmt.Sprintf("%s/anime/%d", BaseURL, malID)

	responseData, err := j.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request %s -> %w", url, err)
	}

	err = json.Unmarshal(responseData, &responseObj)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json data -> %w", err)
	}

	return &responseObj, nil
}
