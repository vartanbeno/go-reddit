package reddit

import (
	"errors"
	"net/http"
	"net/url"
)

// Opt is a configuration option to initialize a client.
type Opt func(*Client) error

// WithHTTPClient sets the HTTP client which will be used to make requests.
func WithHTTPClient(httpClient *http.Client) Opt {
	return func(c *Client) error {
		if httpClient == nil {
			return errors.New("httpClient: cannot be nil")
		}
		c.client = httpClient
		return nil
	}
}

// WithBaseURL sets the base URL for the client to make requests to.
func WithBaseURL(u string) Opt {
	return func(c *Client) error {
		url, err := url.Parse(u)
		if err != nil {
			return err
		}
		c.BaseURL = url
		return nil
	}
}

// WithTokenURL sets the url used to get access tokens.
func WithTokenURL(u string) Opt {
	return func(c *Client) error {
		url, err := url.Parse(u)
		if err != nil {
			return err
		}
		c.TokenURL = url
		return nil
	}
}
