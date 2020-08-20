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

	"github.com/stretchr/testify/require"
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
		"Gold",
		"Listings",
		"Message",
		"Moderation",
		"Multi",
		"Post",
		"Subreddit",
		"User",
	}

	cp := reflect.ValueOf(c)
	cv := reflect.Indirect(cp)

	for _, s := range services {
		require.Falsef(t, cv.FieldByName(s).IsNil(), "c.%s should not be nil", s)
	}
}

func testClientDefaultUserAgent(t *testing.T, c *Client) {
	expectedUserAgent := fmt.Sprintf("golang:%s:v%s (by /u/)", libraryName, libraryVersion)
	require.Equal(t, expectedUserAgent, c.userAgent)
}

func testClientDefaults(t *testing.T, c *Client) {
	testClientDefaultUserAgent(t, c)
	testClientServices(t, c)
}

func TestNewClient(t *testing.T) {
	c, err := NewClient(nil)
	require.NoError(t, err)
	testClientDefaults(t, c)
}

func TestNewClient_Error(t *testing.T) {
	errorOpt := func(c *Client) error {
		return errors.New("foo")
	}

	_, err := NewClient(nil, errorOpt)
	require.EqualError(t, err, "foo")
}

func TestClient_OnRequestComplemented(t *testing.T) {
	setup()
	defer teardown()

	var i int
	cb := func(*http.Request, *http.Response) {
		i++
	}
	client.OnRequestCompleted(cb)

	mux.HandleFunc("/api/v1/test", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
	})

	req, err := client.NewRequest(http.MethodGet, "api/v1/test", nil)
	require.NoError(t, err)

	_, _ = client.Do(ctx, req, nil)
	require.Equal(t, 1, i)

	_, _ = client.Do(ctx, req, nil)
	_, _ = client.Do(ctx, req, nil)
	_, _ = client.Do(ctx, req, nil)
	_, _ = client.Do(ctx, req, nil)
	require.Equal(t, 5, i)

	_, _ = client.Do(ctx, req, nil)
	require.Equal(t, 6, i)
}

func TestClient_JSONErrorResponse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/test", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `{
			"json": {
				"errors": [
					[
						"TEST_ERROR",
						"this is a test error",
						"test field"
					]
				]
			}
		}`)
	})

	req, err := client.NewRequest(http.MethodGet, "api/v1/test", nil)
	require.NoError(t, err)

	resp, err := client.Do(ctx, req, nil)
	require.IsType(t, &JSONErrorResponse{}, err)
	require.EqualError(t, err, fmt.Sprintf(`GET %s/api/v1/test: 200 field "test field" caused TEST_ERROR: this is a test error`, server.URL))
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestClient_ErrorResponse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/test", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{
			"message": "error message"
		}`)
	})

	req, err := client.NewRequest(http.MethodGet, "api/v1/test", nil)
	require.NoError(t, err)

	resp, err := client.Do(ctx, req, nil)
	require.IsType(t, &ErrorResponse{}, err)
	require.EqualError(t, err, fmt.Sprintf(`GET %s/api/v1/test: 403 error message`, server.URL))
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}
