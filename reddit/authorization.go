package reddit

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type basicAuthTransport struct {
	Username string
	Password string
	Base     http.RoundTripper
}

func newBasicAuthTransport(username, password string, base http.RoundTripper) http.RoundTripper {
	return &basicAuthTransport{
		Username: username,
		Password: password,
		Base:     base,
	}
}

func (bat *basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
			bat.Username, bat.Password)))))
	return bat.Base.RoundTrip(req)
}

type accessTokenSource struct {
	Config    *oauth2.Config
	Code      string
	UserAgent string
}

func (s *accessTokenSource) Token() (*oauth2.Token, error) {
	httpClient := &http.Client{Transport: newBasicAuthTransport(s.Config.ClientID, s.Config.ClientSecret,
		&userAgentTransport{
			userAgent: s.UserAgent,
			Base:      &http.Transport{},
		},
	)}
	//httpClient.Transport = helpers.OauthTransport(code, httpClient, lc.Reddit)
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	return s.Config.Exchange(ctx, s.Code)
}
