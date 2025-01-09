package jikan

import (
	"net/url"
	"strconv"
)

type GetTopAnimeQuery struct {
	T      string // type
	Filter string
	Rating string
	SFW    bool
	Page   int
	Limit  int
}

func (q *GetTopAnimeQuery) URL() (string, error) {
	query := url.Values{}
	URL, err := url.Parse(Endpoint + "/top/anime")
	if err != nil {
		return "", err
	}

	if q.Page > 0 {
		query.Add("page", strconv.Itoa(q.Page))
	}

	if q.Limit > 0 {
		query.Add("limit", strconv.Itoa(q.Limit))
	}

	if q.T != "" {
		query.Add("type", q.T)
	}

	URL.RawQuery = query.Encode()

	return URL.String(), nil
}
