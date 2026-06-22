package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Defines all handler function for scionping commands

func (c *Client) StartScionPing(dstURL string, count int) (resp *http.Response, err error) {
	req := PingStartRequest{
		Dst: dstURL,
	}

	fmt.Printf("Ping request to: %s\n", req.Dst)

	// Only set count if > 0
	if count > 0 {
		req.Count = &count
	}

	// Marshal request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Start new ping request
	resp, err = c.client.Post(c.baseURL+ScionPingStartRoute, "application/json", bytes.NewBuffer(body))
	return
}

// Needs doc
func (c *Client) StopScionPing() (resp *http.Response, err error) {
	// Start stop ping request
	resp, err = c.client.Post(c.baseURL+ScionPingStopRoute, "application/json", nil)
	return
}

// Needs doc
func (c *Client) ScionGetResultsPing() (resp *http.Response, err error) {
	resp, err = c.client.Get(c.baseURL + ScionListAvailableRoute)
	return
}

// Needs doc
func (c *Client) ScionStatusPing() (resp *http.Response, err error) {
	resp, err = c.client.Get(c.baseURL + ScionPingStatusRoute)
	return
}
