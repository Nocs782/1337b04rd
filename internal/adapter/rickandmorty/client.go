package rickmorty

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type character struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *Client) FetchCharacterByID(id int) (*character, error) {
	url := fmt.Sprintf("%s/character/%d", c.baseURL, id)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("rickmorty GET error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rickmorty API status: %d", resp.StatusCode)
	}

	var data character
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if data.Image == "" || data.Name == "" {
		return nil, errors.New("character missing name or image")
	}

	return &data, nil
}
