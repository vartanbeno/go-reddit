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

This package currently supports clients for "script" and "web app" options.

2. After creating the app, you will get a client id and client secret.

3. Send a POST request (with the Content-Type header set to "application/x-www-form-urlencoded")
to https://www.reddit.com/api/v1/access_token to obtain the access (and refresh, read further) token.

To authorize a "script" app, include the following form values:
	- grant_type=password
	- username={your Reddit username}
	- password={your Reddit password}

To authorize a "web" app, you will first need to obtain a "code" that can be exchanged for an access token.
It's a two-step process:

3.1 Redirect the user to https://www.reddit.com/api/v1/authorize (Reddit's official authorization URL).
To find out more about the required request parameters, see the "OAuth2" article from the Docs.

3.2. User will be redirected to your app's "Redirect URI". Extract "code" from the query parameters
and exchange it for the access_token, including the following form values:
	- grant_type=authorization_code
	- code={code}
	- redirect_uri={you app's "Redirect URI"}

4. You should receive a response body like the following:
{
	"access_token": "70743860-DRhHVNSEOMu1ldlI",
	"token_type": "bearer",
	"expires_in": 3600,
	"scope": "*"
}

Note: web apps can obtain a refresh token by adding `&duration=permanent` parameter to the "authorization URL" (step 3.1).
*/

package reddit

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// webAppOAuthParams are used to retrieve access token using "code flow" or refresh_token,
// see https://github.com/reddit-archive/reddit/wiki/OAuth2#token-retrieval-code-flow.
type webAppOAuthParams struct {
	// Code can be exchanged for access_token.
	Code string

	// RedirectURI is used to build an AuthCodeURL when requesting users to grant access,
	// and later exchanging code for access_token. The URI must be valid, as it will receive
	// a request containing the `code` after user grants access to the app. Part of the "code flow".
	RedirectURI string

	// RefreshToken should be set to retrieve a new access_token, ignoring the "code flow".
	RefreshToken string
}

// TokenSource creates a reusable token source base on the provided configuration. If code is set,
// it is exchanged for an access_token. If, on the other hand RefreshToken is set, we assume that
// the initial authorization has already happened and create an oauth2.Token with immediate expiry.
func (p webAppOAuthParams) TokenSource(ctx context.Context, config *oauth2.Config) (oauth2.TokenSource, error) {
	var tok *oauth2.Token
	var err error

	if p.RefreshToken != "" {
		tok = &oauth2.Token{
			RefreshToken: p.RefreshToken,
			Expiry:       time.Now(), // refresh before using
		}
	} else if p.Code != "" {
		if tok, err = config.Exchange(ctx, p.Code); err != nil {
			return nil, fmt.Errorf("exchange code: %w", err)
		}
	}
	return config.TokenSource(ctx, tok), err
}

// oauthTokenSource retrieves access_token from resource owner's
// username and password. It implements oauth2.TokenSource.
type oauthTokenSource struct {
	ctx                context.Context
	config             *oauth2.Config
	username, password string
}

func (s *oauthTokenSource) Token() (*oauth2.Token, error) {
	return s.config.PasswordCredentialsToken(s.ctx, s.username, s.password)
}

// oauthTransport returns a Transport to handle authorization based the selected app type.
func oauthTransport(client *Client) (*oauth2.Transport, error) {
	httpClient := &http.Client{Transport: client.client.Transport}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	config := &oauth2.Config{
		ClientID:     client.ID,
		ClientSecret: client.Secret,
		Endpoint: oauth2.Endpoint{
			TokenURL:  client.TokenURL.String(),
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}

	transport := &oauth2.Transport{Base: client.client.Transport}

	switch client.appType {
	case Script:
		transport.Source = oauth2.ReuseTokenSource(nil, &oauthTokenSource{
			ctx:      ctx,
			config:   config,
			username: client.Username,
			password: client.Password,
		})
	case WebApp:
		config.RedirectURL = client.webOauth.RedirectURI
		ts, err := client.webOauth.TokenSource(ctx, config)
		if err != nil {
			return nil, err
		}
		transport.Source = ts
	default:
		// Should we panic here? There is not supposed to be any other app type.
	}

	return transport, nil
}

// AuthCodeURL is a util function for buiding a URL to request permission grant from a user.
//
// TODO: Currently only works with defaultAuthURL,
// but should be able to use a custom AuthURL. Need to find an elegant solution.
//
// By default, Reddit will only issue an access_token to a WebApp for 1h,
// after which the app would need to ask the user to grant access again.
// `permanent` should be set to true to additionally request a refresh_token.
func AuthCodeURL(clientID, redirectURI, state string, scopes []string, permanent bool) string {
	config := &oauth2.Config{
		ClientID: clientID,
		Endpoint: oauth2.Endpoint{
			AuthURL: defaultAuthURL,
		},
		RedirectURL: redirectURI,
		Scopes:      scopes,
	}
	var opts []oauth2.AuthCodeOption
	if permanent {
		opts = append(opts, oauth2.SetAuthURLParam("duration", "permanent"))
	}

	return config.AuthCodeURL(state, opts...)
}
