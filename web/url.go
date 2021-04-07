package web

import (
	"fmt"
	"net/url"

	"github.com/declantraynor/trunc/helpers"
)

type URLBuilder interface {
	Build() *url.URL
}

type RandomURLBuilder struct {
	host   string
	scheme string
}

func (b *RandomURLBuilder) Build() *url.URL {
	return &url.URL{
		Host:   b.host,
		Path:   fmt.Sprintf("/%s", helpers.RandomString(8)),
		Scheme: b.scheme,
	}
}

func NewRandomURLBuilder(baseURL string) (*RandomURLBuilder, error) {
	base, err := helpers.ParseURL(baseURL)
	if err != nil {
		return nil, err
	}
	return &RandomURLBuilder{host: base.Host, scheme: base.Scheme}, nil
}
