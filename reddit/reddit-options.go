package reddit

import (
	"errors"
	"net/http"
	"net/url"
	"os"
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

// FromEnv configures the client with values from environment variables.
// Supported environment variables:
// GO_REDDIT_CLIENT_ID to set the client's id.
// GO_REDDIT_CLIENT_SECRET to set the client's secret.
// GO_REDDIT_CLIENT_USERNAME to set the client's username.
// GO_REDDIT_CLIENT_PASSWORD to set the client's password.
func FromEnv(c *Client) error {
	if v := os.Getenv("GO_REDDIT_CLIENT_ID"); v != "" {
		c.ID = v
	}
	if v := os.Getenv("GO_REDDIT_CLIENT_SECRET"); v != "" {
		c.Secret = v
	}
	if v := os.Getenv("GO_REDDIT_CLIENT_USERNAME"); v != "" {
		c.Username = v
	}
	if v := os.Getenv("GO_REDDIT_CLIENT_PASSWORD"); v != "" {
		c.Password = v
	}
	return nil
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
