package jikan

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
