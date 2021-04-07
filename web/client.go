package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Client struct {
	serverURL string
}

func (c *Client) Shorten(url string) (*ShortenResponse, error) {
	api := fmt.Sprintf("%s/shorten", c.serverURL)
	requestData := &ShortenRequest{URL: url}
	responseData := &ShortenResponse{}

	req := bytes.Buffer{}
	json.NewEncoder(&req).Encode(requestData)

	res, err := http.Post(api, "application/json", &req)
	if err != nil {
		return &ShortenResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		var errText string

		errResponse := &ErrorResponse{}
		if err := json.NewDecoder(res.Body).Decode(&errResponse); err != nil {
			errText = "Invalid response from server"
		} else {
			errText = fmt.Sprintf("HTTP %d response from server: %s", res.StatusCode, errResponse.Error)
		}
		return &ShortenResponse{}, errors.New(errText)
	}

	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		return &ShortenResponse{}, errors.New("Invalid response from server")
	}

	return responseData, nil
}

func NewClient(serverURL string) *Client {
	return &Client{serverURL: serverURL}
}
