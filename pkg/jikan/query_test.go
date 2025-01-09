package jikan_test

import (
	"testing"

	"github.com/reidaa/ano/pkg/jikan"
)

func TestGetTopAnimeQuery(t *testing.T) {
	type test struct {
		input jikan.GetTopAnimeQuery
		want  string
	}

	tests := []test{
		{input: jikan.GetTopAnimeQuery{}, want: "https://api.jikan.moe/v4/top/anime"},
		{input: jikan.GetTopAnimeQuery{Limit: 10}, want: "https://api.jikan.moe/v4/top/anime?limit=10"},
		{input: jikan.GetTopAnimeQuery{Limit: 25, Page: 3}, want: "https://api.jikan.moe/v4/top/anime?limit=25&page=3"},
		{input: jikan.GetTopAnimeQuery{Page: 3, Limit: 3}, want: "https://api.jikan.moe/v4/top/anime?limit=3&page=3"},
		{input: jikan.GetTopAnimeQuery{Limit: 25, Page: 3, T: "tv"}, want: "https://api.jikan.moe/v4/top/anime?limit=25&page=3&type=tv"},
	}

	for _, tc := range tests {
		got, err := tc.input.URL()
		if err != nil {
			t.Fatalf("q.URL() failed with %q", err)
		}
		if got != tc.want {
			t.Errorf("query = %q, want %q", got, tc.want)
		}
	}
}
