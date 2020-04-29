package geddit

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestVoteServiceOp_Up(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("dir", fmt.Sprint(upvote))
		form.Set("rank", "10")

		_ = r.ParseForm()
		if expect, actual := form, r.PostForm; !reflect.DeepEqual(expect, actual) {
			t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
		}

		fmt.Fprint(w, `{}`)
	})

	_, err := client.Vote.Up(ctx, "t3_test")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
}

func TestVoteServiceOp_Down(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("dir", fmt.Sprint(downvote))
		form.Set("rank", "10")

		_ = r.ParseForm()
		if expect, actual := form, r.PostForm; !reflect.DeepEqual(expect, actual) {
			t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
		}

		fmt.Fprint(w, `{}`)
	})

	_, err := client.Vote.Down(ctx, "t3_test")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
}

func TestVoteServiceOp_Remove(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("dir", fmt.Sprint(novote))
		form.Set("rank", "10")

		_ = r.ParseForm()
		if expect, actual := form, r.PostForm; !reflect.DeepEqual(expect, actual) {
			t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
		}

		fmt.Fprint(w, `{}`)
	})

	_, err := client.Vote.Remove(ctx, "t3_test")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
}
