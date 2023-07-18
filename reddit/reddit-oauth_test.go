package reddit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const (
	testCode         = "test_code"
	testAccessToken  = "test_access_token"
	testRefreshToken = "test_refresh_token"
	testRedirectURI  = "http://localhost:5000/auth" // doens't need to be a valid URL

	clientId     = "test_client"
	clientSecret = "test_secret"

	subreddit = "golang"
)

func TestAuthCodeURL(t *testing.T) {
	state := "test_state"
	scopes := []string{"scope_a", "scope_b"}

	for _, tt := range []struct {
		name      string
		permanent bool
	}{
		{"not requesting refresh token", false},
		{"request refresh token", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := url.Parse(AuthCodeURL(clientId, testRedirectURI, state, scopes, tt.permanent))
			if err != nil {
				t.Fatal(err)
			}

			checkQueryParameter(t, got, "client_id", clientId)
			checkQueryParameter(t, got, "state", state)
			checkQueryParameter(t, got, "redirect_uri", testRedirectURI)
			checkQueryParameter(t, got, "scope", strings.Join(scopes, " "))

			if tt.permanent {
				checkQueryParameter(t, got, "duration", "permanent")
			} else {
				checkQueryParameter(t, got, "duration", "")
			}
		})
	}
}

func TestWebAppOauth(t *testing.T) {
	srv := testRedditServer(t)
	t.Cleanup(srv.Close)

	for _, tt := range []struct {
		name string
		opt  Opt
	}{
		{"web app with code", WithWebAppCode(testCode, testRedirectURI)},
		{"web app with refresh_token", WithWebAppRefresh(testRefreshToken)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rc, err := NewClient(
				Credentials{ID: clientId, Secret: clientSecret},
				WithBaseURL(srv.URL),
				WithTokenURL(srv.URL+"/access_token"),
				tt.opt,
			)
			if err != nil {
				t.Fatalf("create client: %v", err)
			}

			// Make a request: check that the client has received the correct access token
			_, _, err = rc.Subreddit.TopPosts(context.Background(), subreddit, nil)
			if err != nil {
				t.Errorf("make authorized request: %v", err)
			}
		})
	}
}

// testRedditServer mocks both reddit.com (for authorization) and oauth.reddit.com (for interactions).
// It only handles a number of endpoints necessary for tests.
func testRedditServer(tb testing.TB) *httptest.Server {
	mux := http.NewServeMux()

	// Exchange code for access_token
	mux.HandleFunc("/access_token", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		w.Header().Set("Content-type", "application/json")

		// Validate grant type
		var ok bool
		switch r.FormValue("grant_type") {
		case "authorization_code":
			if code := r.FormValue("code"); code == testCode {
				// Actual Reddit API returns a different error message
				ok = true
			}
		case "refresh_token":
			if rt := r.FormValue("refresh_token"); rt == testRefreshToken {
				ok = true
			}
		default:
			tb.Log("unexpected grant type:", r.FormValue("grant_type"))
		}

		if !ok {
			// Actual Reddit API returns a different error message
			enc.Encode(map[string]string{"error": "bad_request"})
			return
		}

		enc.Encode(map[string]interface{}{
			"access_token":  testAccessToken,
			"token_type":    "bearer",
			"expires_in":    10 * time.Second,
			"scope":         "scope1,scope2",
			"refresh_token": testRefreshToken,
		})
	})

	// Return the top post for the subreddit
	mux.HandleFunc("/r/"+subreddit+"/top", func(w http.ResponseWriter, r *http.Request) {
		if tok := strings.TrimLeft(r.Header.Get("Authorization"), "Bearer "); tok != testAccessToken {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"kind": kindPost,
			"data": map[string]string{}, // data not needed for the test
		})
	})

	srv := httptest.NewServer(mux)
	return srv
}

// checkQueryParameter validates URL query parameters.
func checkQueryParameter(tb testing.TB, URL *url.URL, param, want string) {
	if got := URL.Query().Get(param); got != want {
		tb.Errorf("%s: got %q, want %q", param, got, want)
	}
}
