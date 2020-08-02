package reddit

import (
	"net/url"
	"os"
)

// Opt is a configuration option to initialize a client.
type Opt func(*Client) error

// FromEnv configures the client with values from environment variables.
//
// Supported environment variables:
// GO_REDDIT_CLIENT_ID to set the client's id.
// GO_REDDIT_CLIENT_SECRET to set the client's secret.
// GO_REDDIT_CLIENT_USERNAME to set the client's username.
// GO_REDDIT_CLIENT_PASSWORD to set the client's password.
func FromEnv(c *Client) error {
	if v, ok := os.LookupEnv("GO_REDDIT_CLIENT_ID"); ok {
		c.ID = v
	}

	if v, ok := os.LookupEnv("GO_REDDIT_CLIENT_SECRET"); ok {
		c.Secret = v
	}

	if v, ok := os.LookupEnv("GO_REDDIT_CLIENT_USERNAME"); ok {
		c.Username = v
	}

	if v, ok := os.LookupEnv("GO_REDDIT_CLIENT_PASSWORD"); ok {
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
