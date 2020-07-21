package reddit

import (
	"os"
	"reflect"
	"testing"
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
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	type values struct {
		id, secret, username, password string
	}

	expect := values{"id1", "secret1", "username1", "password1"}
	actual := values{c.ID, c.Secret, c.Username, c.Password}

	if !reflect.DeepEqual(expect, actual) {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}
}

func TestWithBaseURL(t *testing.T) {
	baseURL := "http://localhost:8080"
	c, err := NewClient(nil, WithBaseURL(baseURL))
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	if expect, actual := baseURL, c.BaseURL.String(); expect != actual {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}
}

func TestWithTokenURL(t *testing.T) {
	tokenURL := "http://localhost:8080/api/v1/access_token"
	c, err := NewClient(nil, WithTokenURL(tokenURL))
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	if expect, actual := tokenURL, c.TokenURL.String(); expect != actual {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}
}
