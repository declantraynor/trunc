package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInvalidURL(t *testing.T) {
	urls := []string{
		"http://",
		"localhost",
		"localhost:5000",
		"://",
	}
	for _, u := range urls {
		parsed, err := ParseURL(u)
		assert.Nil(t, parsed)
		assert.NotNil(t, err)
	}
}

func TestParseURL(t *testing.T) {
	parsed, err := ParseURL("https://example.com/path")
	assert.Equal(t, parsed.Scheme, "https")
	assert.Equal(t, parsed.Host, "example.com")
	assert.Equal(t, parsed.Path, "/path")
	assert.Nil(t, err)
}
