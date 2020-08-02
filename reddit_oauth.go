/*
Docs:
- https://www.reddit.com/dev/api/
- https://github.com/reddit-archive/reddit/wiki/api
- https://github.com/reddit-archive/reddit/wiki/OAuth2
- https://github.com/reddit-archive/reddit/wiki/OAuth2-Quick-Start-Example

1. Go to https://www.reddit.com/prefs/apps and create an app. There are 3 types of apps:
	- Web app. Service is available over http or https, preferably the latter.
	- Installed app, such as a mobile app on a user's device which you can't control.
	  Redirect the user to a URI after they grant your app permissions.
	- Script (the simplest type of app). Select this if you are the only person who will
	  use the app. Only has access to your account.

Best option for a client like this is to use the script option.

2. After creating the app, you will get a client id and client secret.

3. Send a POST request (with the Content-Type header set to "application/x-www-form-urlencoded")
to https://www.reddit.com/api/v1/access_token with the following form values:
	- grant_type=password
	- username={your Reddit username}
	- password={your Reddit password}

4. You should receive a response body like the following:
{
	"access_token": "70743860-DRhHVNSEOMu1ldlI",
	"token_type": "bearer",
	"expires_in": 3600,
	"scope": "*"
}
*/

package reddit

import (
	"context"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var endpoint = oauth2.Endpoint{
	TokenURL:  "https://www.reddit.com/api/v1/access_token",
	AuthStyle: oauth2.AuthStyleInHeader,
}

type oauth2Config struct {
	id       string
	secret   string
	username string
	password string
	tokenURL string

	// We need to set a custom user agent, because using the one set by default by the
	// stdlib gives us 429 Too Many Request responses from the Reddit API.
	userAgentTransport *userAgentTransport
}

func oauth2Transport(c oauth2Config) *oauth2.Transport {
	params := url.Values{
		"grant_type": {"password"},
		"username":   {c.username},
		"password":   {c.password},
	}

	cfg := clientcredentials.Config{
		ClientID:       c.id,
		ClientSecret:   c.secret,
		TokenURL:       c.tokenURL,
		AuthStyle:      oauth2.AuthStyleInHeader,
		EndpointParams: params,
	}

	httpClient := &http.Client{Transport: c.userAgentTransport}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	src := cfg.TokenSource(ctx)
	tr := &oauth2.Transport{
		Source: src,
		Base:   c.userAgentTransport,
	}
	return tr
}

// WithCredentials sets the necessary values for the client to authenticate via OAuth2.
func WithCredentials(id, secret, username, password string) Opt {
	return func(c *Client) error {
		c.ID = id
		c.Secret = secret
		c.Username = username
		c.Password = password
		return nil
	}
}
