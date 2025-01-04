package jikan

type Title struct {
	Type  string `json:"type"`
	Title string `json:"title"`
}

type Image struct {
	ImageURL      string `json:"image_url"`
	SmallImageURL string `json:"small_image_url"`
	LargeImageURL string `json:"large_image_url"`
}

type AnimeImage struct {
	Jpg  Image `json:"jpg"`
	Webp Image `json:"webp"`
}

type Anime struct {
	Images     AnimeImage `json:"images"`
	URL        string     `json:"url"`
	Type       string     `json:"type"`
	Status     string     `json:"status"`
	Titles     []Title    `json:"titles"`
	MalID      int        `json:"mal_id"`
	ScoredBy   int        `json:"scored_by"`
	Rank       int        `json:"rank"` // Ranking are not accurate
	Popularity int        `json:"popularity"`
	Members    int        `json:"members"`
	Favorites  int        `json:"favorites"`
	Score      float32    `json:"score"`
}

type Item struct {
	Count   int `json:"count"`
	Total   int `json:"total"`
	PerPage int `json:"per_page"`
}

type Pagination struct {
	LastVisiblePage int  `json:"last_visible_page"`
	HasNextPage     bool `json:"has_next_page"`
	Items           Item `json:"items"`
}
type TopAnimeResponse struct {
	Data       []Anime    `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type AnimeResponse struct {
	Data Anime `json:"data"`
}

type ErrorResponse struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	Error     string `json:"error"`
	ReportUrl string `json:"report_url"`
	Status    int    `json:"status"`
}
