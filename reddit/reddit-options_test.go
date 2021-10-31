package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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
	require.Equal(t, "id1", c.Oauth2Config.ClientID)
	require.Equal(t, "secret1", c.Oauth2Config.ClientSecret)
	require.Equal(t, "username1", c.Username)
	require.Equal(t, "password1", c.Password)
}
