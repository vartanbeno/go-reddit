package reddit

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithHTTPClient(t *testing.T) {
	_, err := NewClient(&Credentials{}, WithHTTPClient(nil))
	require.EqualError(t, err, "httpClient: cannot be nil")

	_, err = NewClient(&Credentials{}, WithHTTPClient(&http.Client{}))
	require.NoError(t, err)
}

func TestWithBaseURL(t *testing.T) {
	c, err := NewClient(&Credentials{}, WithBaseURL(":"))
	urlErr, ok := err.(*url.Error)
	require.True(t, ok)
	require.Equal(t, "parse", urlErr.Op)

	baseURL := "http://localhost:8080"
	c, err = NewClient(&Credentials{}, WithBaseURL(baseURL))
	require.NoError(t, err)
	require.Equal(t, baseURL, c.BaseURL.String())
}

func TestWithTokenURL(t *testing.T) {
	c, err := NewClient(&Credentials{}, WithTokenURL(":"))
	urlErr, ok := err.(*url.Error)
	require.True(t, ok)
	require.Equal(t, "parse", urlErr.Op)

	tokenURL := "http://localhost:8080/api/v1/access_token"
	c, err = NewClient(&Credentials{}, WithTokenURL(tokenURL))
	require.NoError(t, err)
	require.Equal(t, tokenURL, c.TokenURL.String())
}
