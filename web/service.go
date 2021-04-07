package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/declantraynor/trunc/helpers"
	"github.com/declantraynor/trunc/storage"
)

type Service struct {
	Store      storage.Store
	URLBuilder URLBuilder
}

// Shorten handles requests which supply a valid URL for shortening.
// The "shortening" is implemented by generating a new URL whose path
// is used as a key for storing the original URL in a backing store.
//
// The new URL is then returned to the client. If the client happens
// to visit that URL, this service can look up the original URL and
// redirect accordingly.
func (s *Service) Shorten(w http.ResponseWriter, r *http.Request) {
	req, err := parseShortenRequest(r)
	if err != nil {
		renderError(w, http.StatusBadRequest, err)
		return
	}

	shortURL := s.URLBuilder.Build()
	if err := s.Store.Set(shortURL.Path, req.URL); err != nil {
		renderError(w, http.StatusInternalServerError, nil)
		return
	}

	response := ShortenResponse{
		LongURL:  req.URL,
		ShortURL: shortURL.String(),
	}
	renderJSON(w, response, http.StatusOK)
}

// Redirect assumes that the path portion of the incoming request URL
// has been generated by this service, in which case it can be used to
// retrieve a previously stored ("shortened") URL.
//
// If a stored URL is found, this handler redirects the client to the
// stored URL. If not, the handler responds with a HTTP 404.
func (s *Service) Redirect(w http.ResponseWriter, r *http.Request) {
	longURL, err := s.Store.Get(r.URL.Path)
	if err != nil {
		renderError(w, http.StatusNotFound, nil)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

func parseShortenRequest(r *http.Request) (*ShortenRequest, error) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	parsed := &ShortenRequest{}
	if err := json.Unmarshal(body, parsed); err != nil {
		return &ShortenRequest{}, errors.New("Request data must be valid JSON")
	}

	if parsed.URL == "" {
		return &ShortenRequest{}, errors.New("Missing required data: url")
	}

	if _, err := helpers.ParseURL(parsed.URL); err != nil {
		return &ShortenRequest{}, errors.New(fmt.Sprintf("Invalid data: [url]=%s", parsed.URL))
	}

	return parsed, nil
}

func renderError(w http.ResponseWriter, status int, err error) {
	var description string
	if err == nil {
		description = http.StatusText(status)
	} else {
		description = err.Error()
	}
	response := ErrorResponse{Error: description}
	renderJSON(w, response, status)
}

func renderJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
