package reddit

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestWithHTTPClient(t *testing.T) {
	_, err := NewClient(Credentials{}, WithHTTPClient(nil))
	require.EqualError(t, err, "*http.Client: cannot be nil")

	_, err = NewClient(Credentials{}, WithHTTPClient(&http.Client{}))
	require.NoError(t, err)
}

func TestWithUserAgent(t *testing.T) {
	c, err := NewClient(Credentials{}, WithUserAgent("test"))
	require.NoError(t, err)
	require.Equal(t, "test", c.UserAgent())

	c, err = NewClient(Credentials{}, WithUserAgent(""))
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("golang:%s:v%s", libraryName, libraryVersion), c.UserAgent())
}

func TestWithBaseURL(t *testing.T) {
	c, err := NewClient(Credentials{}, WithBaseURL(":"))
	urlErr, ok := err.(*url.Error)
	require.True(t, ok)
	require.Equal(t, "parse", urlErr.Op)

	baseURL := "http://localhost:8080"
	c, err = NewClient(Credentials{}, WithBaseURL(baseURL))
	require.NoError(t, err)
	require.Equal(t, baseURL, c.BaseURL.String())
}

func TestWithTokenURL(t *testing.T) {
	c, err := NewClient(Credentials{}, WithTokenURL(":"))
	urlErr, ok := err.(*url.Error)
	require.True(t, ok)
	require.Equal(t, "parse", urlErr.Op)

	tokenURL := "http://localhost:8080/api/v1/access_token"
	c, err = NewClient(Credentials{}, WithTokenURL(tokenURL))
	require.NoError(t, err)
	require.Equal(t, tokenURL, c.TokenURL.String())
}

type RequestInterceptor struct {
	interceptedBody string
}

func (t *RequestInterceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	t.interceptedBody = string(requestBody)
	var body bytes.Buffer
	body.WriteString(`{"access_token": "foobar", "expires_in": 3600, "scope": "*", "token_type": "bearer"}`)
	return &http.Response{Status: "200 OK", StatusCode: 200, Body: io.NopCloser(&body)}, nil
}

func TestWithApplicationOnlyOAuth(t *testing.T) {
	requestInterceptor := &RequestInterceptor{}
	c, err := NewClient(Credentials{ID: "id", Secret: "secret"}, WithApplicationOnlyOAuth(true), WithHTTPClient(&http.Client{Transport: requestInterceptor}))
	require.NoError(t, err)
	token, err := c.client.Transport.(*oauth2.Transport).Source.Token()
	require.NoError(t, err)
	require.Equal(t, token.AccessToken, "foobar")
	require.Equal(t, "grant_type=client_credentials", requestInterceptor.interceptedBody)
}

func TestFromEnv(t *testing.T) {
	os.Setenv("GO_REDDIT_CLIENT_ID", "id1")
	defer os.Unsetenv("GO_REDDIT_CLIENT_ID")

	os.Setenv("GO_REDDIT_CLIENT_SECRET", "secret1")
	defer os.Unsetenv("GO_REDDIT_CLIENT_SECRET")

	os.Setenv("GO_REDDIT_CLIENT_USERNAME", "username1")
	defer os.Unsetenv("GO_REDDIT_CLIENT_USERNAME")

	os.Setenv("GO_REDDIT_CLIENT_PASSWORD", "password1")
	defer os.Unsetenv("GO_REDDIT_CLIENT_PASSWORD")

	c, err := NewClient(Credentials{}, FromEnv)
	require.NoError(t, err)
	require.Equal(t, "id1", c.ID)
	require.Equal(t, "secret1", c.Secret)
	require.Equal(t, "username1", c.Username)
	require.Equal(t, "password1", c.Password)
}
