package bhr

import (
	"log"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: time.Second * 10},
		apiKey:     apiKey,
	}
}

func (c *Client) Request(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(c.apiKey, "")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-agent", "bhr/0.0.1")
	return c.httpClient.Do(req)
}
