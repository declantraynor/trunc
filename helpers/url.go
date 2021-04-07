package helpers

import (
	"errors"
	"net/url"
)

func ParseURL(u string) (*url.URL, error) {
	parsed, err := url.ParseRequestURI(u)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme == "" {
		return nil, errors.New("URL must have a scheme (e.g. http://)")
	}
	if parsed.Host == "" {
		return nil, errors.New("URL must have a host (e.g. example.com)")
	}

	return parsed, nil
}
