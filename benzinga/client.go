package benzinga

import "net/http"

const (
	DefaultHostname = "api.benzinga.com"
)

type Client struct {
	c *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		c: httpClient,
	}
}
