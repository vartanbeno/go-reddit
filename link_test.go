package geddit

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestLinkServiceOp_Hide(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/hide", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		form := url.Values{}
		form.Set("id", "1,2,3")

		_ = r.ParseForm()
		if expect, actual := form, r.PostForm; !reflect.DeepEqual(expect, actual) {
			t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
		}

		fmt.Fprint(w, `{}`)
	})

	_, err := client.Link.Hide(ctx)
	if err == nil {
		t.Fatal("expected error, got nothing instead")
	}
	if expect, actual := `must provide at least 1 id`, err.Error(); expect != actual {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}

	res, err := client.Link.Hide(ctx, "1", "2", "3")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	if expect, actual := http.StatusOK, res.StatusCode; expect != actual {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}
}

func TestLinkServiceOp_Unhide(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unhide", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		form := url.Values{}
		form.Set("id", "1,2,3")

		_ = r.ParseForm()
		if expect, actual := form, r.PostForm; !reflect.DeepEqual(expect, actual) {
			t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
		}

		fmt.Fprint(w, `{}`)
	})

	_, err := client.Link.Unhide(ctx)
	if err == nil {
		t.Fatal("expected error, got nothing instead")
	}
	if expect, actual := `must provide at least 1 id`, err.Error(); expect != actual {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}

	res, err := client.Link.Unhide(ctx, "1", "2", "3")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	if expect, actual := http.StatusOK, res.StatusCode; expect != actual {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}
}
