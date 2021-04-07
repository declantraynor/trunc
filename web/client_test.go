package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockHandler(status int, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(data)
	}
}

func TestClientBadRequest(t *testing.T) {
	mockResponse := &ErrorResponse{Error: "Bad Request"}
	server := httptest.NewServer(http.HandlerFunc(mockHandler(400, mockResponse)))
	defer server.Close()

	client := NewClient(server.URL)
	response, err := client.Shorten("http://example.com")

	assert.Empty(t, response)
	assert.Equal(t, "HTTP 400 response from server: Bad Request", err.Error())
}

func TestClient(t *testing.T) {
	mockResponse := &ShortenResponse{LongURL: "http://example.com", ShortURL: "http://trunc/b2erAs13"}
	server := httptest.NewServer(http.HandlerFunc(mockHandler(200, mockResponse)))
	defer server.Close()

	client := NewClient(server.URL)
	response, err := client.Shorten("http://example.com")

	assert.Equal(t, mockResponse, response)
	assert.Nil(t, err)
}
