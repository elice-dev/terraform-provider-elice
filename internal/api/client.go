package api

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	HostURL      string
	HTTPClient   *http.Client
	Token        string
	Organization string
}

func NewClient(host, token, organization *string) (*Client, error) {
	c := Client{
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
		HostURL:      *host,
		Token:        *token,
		Organization: *organization,
	}

	if host != nil {
		c.HostURL = *host
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("x-elice-org-name-short", c.Organization)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
