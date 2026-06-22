package api

import (
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

type ClientConfig struct {
	BaseURL string
	Timeout time.Duration
	// Maybe some default headers?
}

func NewClient(clientConfig ClientConfig) *Client {
	return &Client{
		baseURL: clientConfig.BaseURL,
		client: &http.Client{
			Timeout: clientConfig.Timeout,
		},
	}
}
