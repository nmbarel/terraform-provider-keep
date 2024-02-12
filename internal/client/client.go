package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client struct with Api Key needed to authenticate against keep
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	ApiKey     string
}

// NewClient func creates new client 
func NewClient(host_url string, api_key string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL: host_url,
		ApiKey: api_key,
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Add("X-API-KEY", c.ApiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	statusOk := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOk{
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}