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
	"time"

	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func setup(t testing.TB) (*Client, *http.ServeMux) {
	mux := http.NewServeMux()

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	mux.HandleFunc("/api/v1/access_token", func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"access_token": "token1",
			"token_type": "bearer",
			"expires_in": 3600,
			"scope": "*"
		}`
		w.Header().Add(headerContentType, mediaTypeJSON)
		fmt.Fprint(w, response)
	})

	client, _ := NewClient(
		Credentials{"id1", "secret1", "user1", "password1"},
		WithBaseURL(server.URL),
		WithTokenURL(server.URL+"/api/v1/access_token"),
	)

	return client, mux
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
		"LiveThread",
		"Message",
		"Moderation",
		"Multi",
		"Post",
		"Stream",
		"Subreddit",
		"User",
		"Widget",
		"Wiki",
	}

	cp := reflect.ValueOf(c)
	cv := reflect.Indirect(cp)

	for _, s := range services {
		require.Falsef(t, cv.FieldByName(s).IsNil(), "c.%s should not be nil", s)
	}
}

func testClientDefaultUserAgent(t *testing.T, c *Client) {
	expectedUserAgent := fmt.Sprintf("golang:%s:v%s", libraryName, libraryVersion)
	require.Equal(t, expectedUserAgent, c.UserAgent())
}

func testClientDefaults(t *testing.T, c *Client) {
	testClientDefaultUserAgent(t, c)
	testClientServices(t, c)
}

func TestNewClient(t *testing.T) {
	c, err := NewClient(Credentials{})
	require.NoError(t, err)
	testClientDefaults(t, c)
}

func TestNewClient_Error(t *testing.T) {
	_, err := NewClient(Credentials{})
	require.NoError(t, err)

	errorOpt := func(c *Client) error {
		return errors.New("foo")
	}

	_, err = NewClient(Credentials{}, errorOpt)
	require.EqualError(t, err, "foo")
}

func TestNewReadonlyClient(t *testing.T) {
	c, err := NewReadonlyClient()
	require.NoError(t, err)
	require.Equal(t, c.BaseURL.String(), defaultBaseURLReadonly)
}

func TestNewReadonlyClient_Error(t *testing.T) {
	_, err := NewReadonlyClient()
	require.NoError(t, err)

	errorOpt := func(c *Client) error {
		return errors.New("foo")
	}

	_, err = NewReadonlyClient(errorOpt)
	require.EqualError(t, err, "foo")
}

func TestDefaultClient(t *testing.T) {
	require.NotNil(t, DefaultClient())
}

func TestClient_Readonly_NewRequest(t *testing.T) {
	c, err := NewReadonlyClient()
	require.NoError(t, err)

	req, err := c.NewRequest(http.MethodGet, "r/golang", nil)
	require.NoError(t, err)
	require.Equal(t, defaultBaseURLReadonly+"/r/golang.json", req.URL.String())
}

func TestClient_OnRequestComplemented(t *testing.T) {
	client, mux := setup(t)

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
	client, mux := setup(t)

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
	require.EqualError(t, err, fmt.Sprintf(`GET %s/api/v1/test: 200 field "test field" caused TEST_ERROR: this is a test error`, client.BaseURL))
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestClient_ErrorResponse(t *testing.T) {
	client, mux := setup(t)

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
	require.EqualError(t, err, fmt.Sprintf(`GET %s/api/v1/test: 403 error message`, client.BaseURL))
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestClient_Do_RateLimitError(t *testing.T) {
	client, mux := setup(t)

	var counter int
	mux.HandleFunc("/api/v1/test", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		defer func() { counter++ }()

		switch counter {
		case 0:
			w.Header().Set(headerRateLimitRemaining, "500")
			w.Header().Set(headerRateLimitUsed, "100")
			w.Header().Set(headerRateLimitReset, "120")
		case 1:
			w.Header().Set(headerRateLimitRemaining, "0")
			w.Header().Set(headerRateLimitUsed, "600")
			w.Header().Set(headerRateLimitReset, "240")
		}
	})

	req, err := client.NewRequest(http.MethodGet, "api/v1/test", nil)
	require.NoError(t, err)

	client.rate.Remaining = 0
	client.rate.Reset = time.Now().Add(time.Minute)

	resp, err := client.Do(ctx, req, nil)
	require.Equal(t, 0, counter)
	require.IsType(t, &RateLimitError{}, err)
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	client.rate = Rate{}

	resp, err = client.Do(ctx, req, nil)
	require.Equal(t, 1, counter)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 500, resp.Rate.Remaining)
	require.Equal(t, 100, resp.Rate.Used)
	require.Equal(t, time.Now().Truncate(time.Second).Add(time.Minute*2), resp.Rate.Reset)

	resp, err = client.Do(ctx, req, nil)
	require.Equal(t, 2, counter)
	require.IsType(t, &RateLimitError{}, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 0, resp.Rate.Remaining)
	require.Equal(t, 600, resp.Rate.Used)
	require.Equal(t, time.Now().Truncate(time.Second).Add(time.Minute*4), resp.Rate.Reset)
}
