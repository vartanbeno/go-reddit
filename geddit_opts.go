package geddit

import "fmt"

// Opt is a configuration option to initialize a client
type Opt func(*Client) error

// WithUserAgent sets the user agent for the client
func WithUserAgent(ua string) Opt {
	return func(c *Client) error {
		c.UserAgent = fmt.Sprintf("%s %s", ua, c.UserAgent)
		return nil
	}
}
