package reddit

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromEnv(t *testing.T) {
	os.Setenv("GO_REDDIT_CLIENT_ID", "id1")
	defer os.Unsetenv("GO_REDDIT_CLIENT_ID")

	os.Setenv("GO_REDDIT_CLIENT_SECRET", "secret1")
	defer os.Unsetenv("GO_REDDIT_CLIENT_SECRET")

	os.Setenv("GO_REDDIT_CLIENT_USERNAME", "username1")
	defer os.Unsetenv("GO_REDDIT_CLIENT_USERNAME")

	os.Setenv("GO_REDDIT_CLIENT_PASSWORD", "password1")
	defer os.Unsetenv("GO_REDDIT_CLIENT_PASSWORD")

	c, err := NewClient(nil, FromEnv)
	require.NoError(t, err)

	type values struct {
		id, secret, username, password string
	}

	expect := values{"id1", "secret1", "username1", "password1"}
	actual := values{c.ID, c.Secret, c.Username, c.Password}
	require.Equal(t, expect, actual)
}

func TestWithCredentials(t *testing.T) {
	withCredentials := WithCredentials("id1", "secret1", "username1", "password1")
	c, err := NewClient(nil, withCredentials)
	require.NoError(t, err)
	require.Equal(t, "id1", c.ID)
	require.Equal(t, "secret1", c.Secret)
	require.Equal(t, "username1", c.Username)
	require.Equal(t, "password1", c.Password)
}

func TestWithBaseURL(t *testing.T) {
	baseURL := "http://localhost:8080"
	c, err := NewClient(nil, WithBaseURL(baseURL))
	require.NoError(t, err)
	require.Equal(t, baseURL, c.BaseURL.String())
}

func TestWithTokenURL(t *testing.T) {
	tokenURL := "http://localhost:8080/api/v1/access_token"
	c, err := NewClient(nil, WithTokenURL(tokenURL))
	require.NoError(t, err)
	require.Equal(t, tokenURL, c.TokenURL.String())
}
