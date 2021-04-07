package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockStore) Set(key, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

type MockURLBuilder struct{}

func (m *MockURLBuilder) Build() *url.URL {
	return &url.URL{
		Host:   "trunc",
		Path:   "/cTRD2S4d",
		Scheme: "http",
	}
}

var mockStore = &MockStore{}
var service = &Service{
	Store:      mockStore,
	URLBuilder: &MockURLBuilder{},
}

func TestShorten(t *testing.T) {
	mockStore.
		On(
			"Set",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
		Return(nil).
		Once()

	body := strings.NewReader(`{"url": "http://example.com/really/long/url?foo=bar123456789"}`)
	request, _ := http.NewRequest("POST", "http://trunc/shorten", body)
	expectJSON := map[string]string{
		"long_url":  "http://example.com/really/long/url?foo=bar123456789",
		"short_url": "http://trunc/cTRD2S4d",
	}
	assertResponse(t, service.Shorten, request, http.StatusOK, expectJSON)

	mockStore.AssertExpectations(t)
}

func TestShortenStorageError(t *testing.T) {
	mockStore.
		On(
			"Set",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
		Return(errors.New("storage error")).
		Once()

	body := strings.NewReader(`{"url": "http://example.com/really/long/url?foo=bar123456789"}`)
	request, _ := http.NewRequest("POST", "http://trunc/shorten", body)
	expectJSON := map[string]string{
		"error": http.StatusText(http.StatusInternalServerError),
	}
	assertResponse(t, service.Shorten, request, http.StatusInternalServerError, expectJSON)

	mockStore.AssertExpectations(t)
}

func TestShortenInvalidJSON(t *testing.T) {
	body := strings.NewReader(`not json`)
	request, _ := http.NewRequest("POST", "http://trunc/shorten", body)
	expectJSON := map[string]string{
		"error": "Request data must be valid JSON",
	}
	assertResponse(t, service.Shorten, request, http.StatusBadRequest, expectJSON)
}

func TestShortenMissingURL(t *testing.T) {
	body := strings.NewReader(`{}`)
	request, _ := http.NewRequest("POST", "http://trunc/shorten", body)
	expectJSON := map[string]string{
		"error": "Missing required data: url",
	}
	assertResponse(t, service.Shorten, request, http.StatusBadRequest, expectJSON)
}

func TestShortenInvalidURL(t *testing.T) {
	body := strings.NewReader(`{"url": "not_a_url"}`)
	request, _ := http.NewRequest("POST", "http://trunc/shorten", body)
	expectJSON := map[string]string{
		"error": "Invalid data: [url]=not_a_url",
	}
	assertResponse(t, service.Shorten, request, http.StatusBadRequest, expectJSON)
}

func TestRedirect(t *testing.T) {
	expectRedirectURL := "http://example.com/really/long/url?foo=bar123456789"

	mockStore.
		On(
			"Get",
			"/cTRD2S4d",
		).
		Return(expectRedirectURL, nil).
		Once()

	request := httptest.NewRequest("GET", "http://trunc/cTRD2S4d", nil)
	response := assertResponse(t, service.Redirect, request, http.StatusMovedPermanently, nil)

	assert.Equal(t, response.Header().Get("Location"), expectRedirectURL)

	mockStore.AssertExpectations(t)
}

func TestRedirectNotFound(t *testing.T) {
	mockStore.
		On(
			"Get",
			"/cTRD2S4d",
		).
		Return("", errors.New("key not found")).
		Once()

	request := httptest.NewRequest("GET", "http://trunc/cTRD2S4d", nil)
	assertResponse(t, service.Redirect, request, http.StatusNotFound, nil)

	mockStore.AssertExpectations(t)
}

// assertResponse is a convenience intended to reduce duplication in the test
// suite. Since most handler tests do the same things, this function can take
// a reference to a test, along with references to the handler under test, its
// inputs and expected outputs, and perform some assertions.
//
// The generated response is returned so that it can be further asserted on,
// if required.
func assertResponse(
	t *testing.T,
	handler http.HandlerFunc,
	request *http.Request,
	expectCode int,
	expectJSON map[string]string,
) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()

	handler(response, request)
	assert.Equal(t, response.Code, expectCode)

	if expectJSON != nil {
		var b bytes.Buffer
		json.NewEncoder(&b).Encode(expectJSON)
		assert.Equal(t, b.String(), response.Body.String())
	}

	return response
}
