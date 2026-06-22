package api

import "net/http"

// Defines all handler function for capture commands

func (c *Client) StartCapture() (resp *http.Response, err error) {
	// Leaving out interface for now TODO - sending no body
	resp, err = c.client.Post(c.baseURL+CaptureStartRoute, "application/json", nil)
	return
}

func (c *Client) StopCapture() (resp *http.Response, err error) {
	// Leaving out interface for now TODO - sending no body
	resp, err = c.client.Post(c.baseURL+CaptureStopRoute, "application/json", nil)
	return
}

func (c *Client) StatusCapture() (resp *http.Response, err error) {
	// Leaving out interface for now TODO - sending no body
	resp, err = c.client.Get(c.baseURL + CaptureStatusRoute)
	return
}

func (c *Client) GetResultsCapture() (resp *http.Response, err error) {
	resp, err = c.client.Get(c.baseURL + CaptureListAvailableRoute)
	return
}
