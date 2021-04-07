package web

type ShortenResponse struct {
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
