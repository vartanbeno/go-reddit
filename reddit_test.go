package reddit

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mux    *http.ServeMux
	ctx    = context.Background()
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	mux.HandleFunc("/api/v1/access_token", func(w http.ResponseWriter, r *http.Request) {
		response := `
		{
			"access_token": "token1",
			"token_type": "bearer",
			"expires_in": 3600,
			"scope": "*"
		}
		`
		w.Header().Add(headerContentType, mediaTypeJSON)
		fmt.Fprint(w, response)
	})

	client, _ = NewClient(nil,
		WithCredentials("id1", "secret1", "user1", "password1"),
		WithBaseURL(server.URL),
		WithTokenURL(server.URL+"/api/v1/access_token"),
	)
}

func teardown() {
	server.Close()
}

func readFileContents(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func testClientServices(t *testing.T, c *Client) {
	services := []string{
		"Account",
		"Collection",
		"Comment",
		"Emoji",
		"Flair",
		"Listings",
		"Moderation",
		"Multi",
		"Post",
		"Subreddit",
		"User",
	}

	cp := reflect.ValueOf(c)
	cv := reflect.Indirect(cp)

	for _, s := range services {
		assert.Falsef(t, cv.FieldByName(s).IsNil(), "c.%s should not be nil", s)
	}
}

func testClientDefaultUserAgent(t *testing.T, c *Client) {
	expectedUserAgent := fmt.Sprintf("golang:%s:v%s (by /u/)", libraryName, libraryVersion)
	assert.Equal(t, expectedUserAgent, c.userAgent)
}

func testClientDefaults(t *testing.T, c *Client) {
	testClientDefaultUserAgent(t, c)
	testClientServices(t, c)
}

func TestNewClient(t *testing.T) {
	c, err := NewClient(nil)
	assert.NoError(t, err)
	testClientDefaults(t, c)
}

func TestNewClient_Error(t *testing.T) {
	errorOpt := func(c *Client) error {
		return errors.New("foo")
	}

	_, err := NewClient(nil, errorOpt)
	assert.EqualError(t, err, "foo")
}
