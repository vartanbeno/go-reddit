package geddit

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

var expectedCommentSubmitOrEdit = &Comment{
	ID:        "test2",
	FullID:    "t1_test2",
	ParentID:  "t1_test",
	Permalink: "/r/subreddit/comments/test1/some_thread/test2/",

	Body:            "test comment",
	BodyHTML:        "<div class=\"md\"><p>test comment</p>\n</div>",
	Author:          "reddit_username",
	AuthorID:        "t2_user1",
	AuthorFlairText: "Flair",

	Subreddit:             "subreddit",
	SubredditNamePrefixed: "r/subreddit",
	SubredditID:           "t5_test",

	Score:            1,
	Controversiality: 0,

	Created:    1588147787,
	CreatedUTC: 1588118987,

	LinkID: "t3_link1",
}

func TestCommentServiceOp_Submit(t *testing.T) {
	setup()
	defer teardown()

	commentBlob := readFileContents(t, "testdata/comment-submit-edit.json")

	mux.HandleFunc("/api/comment", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("return_rtjson", "true")
		form.Set("parent", "t1_test")
		form.Set("text", "test comment")

		_ = r.ParseForm()
		if expect, actual := form, r.PostForm; !reflect.DeepEqual(expect, actual) {
			t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
		}

		fmt.Fprint(w, commentBlob)
	})

	comment, _, err := client.Comment.Submit(ctx, "t1_test", "test comment")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	if expect, actual := expectedCommentSubmitOrEdit, comment; !reflect.DeepEqual(expect, actual) {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}
}

func TestCommentServiceOp_Edit(t *testing.T) {
	setup()
	defer teardown()

	commentBlob := readFileContents(t, "testdata/comment-submit-edit.json")

	mux.HandleFunc("/api/editusertext", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("return_rtjson", "true")
		form.Set("thing_id", "t1_test")
		form.Set("text", "test comment")

		_ = r.ParseForm()
		if expect, actual := form, r.PostForm; !reflect.DeepEqual(expect, actual) {
			t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
		}

		fmt.Fprint(w, commentBlob)
	})

	_, _, err := client.Comment.Edit(ctx, "t3_test", "test comment")
	if err == nil {
		t.Fatal("expected error, got nothing instead")
	}
	if expect, actual := `must provide comment id (starting with t1_); id provided: "t3_test"`, err.Error(); expect != actual {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}

	comment, _, err := client.Comment.Edit(ctx, "t1_test", "test comment")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	if expect, actual := expectedCommentSubmitOrEdit, comment; !reflect.DeepEqual(expect, actual) {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}
}

func TestCommentServiceOp_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/del", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		form := url.Values{}
		form.Set("id", "t1_test")

		_ = r.ParseForm()
		if expect, actual := form, r.PostForm; !reflect.DeepEqual(expect, actual) {
			t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
		}

		fmt.Fprint(w, `{}`)
	})

	_, err := client.Comment.Delete(ctx, "t3_test")
	if err == nil {
		t.Fatal("expected error, got nothing instead")
	}
	if expect, actual := `must provide comment id (starting with t1_); id provided: "t3_test"`, err.Error(); expect != actual {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}

	res, err := client.Comment.Delete(ctx, "t1_test")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	if expect, actual := http.StatusOK, res.StatusCode; expect != actual {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}
}

func TestCommentServiceOp_Save(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/save", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		form := url.Values{}
		form.Set("id", "t1_test")

		_ = r.ParseForm()
		if expect, actual := form, r.PostForm; !reflect.DeepEqual(expect, actual) {
			t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
		}

		fmt.Fprint(w, `{}`)
	})

	_, err := client.Comment.Save(ctx, "t3_test")
	if err == nil {
		t.Fatal("expected error, got nothing instead")
	}
	if expect, actual := `must provide comment id (starting with t1_); id provided: "t3_test"`, err.Error(); expect != actual {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}

	res, err := client.Comment.Save(ctx, "t1_test")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	if expect, actual := http.StatusOK, res.StatusCode; expect != actual {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}
}

func TestCommentServiceOp_Unsave(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unsave", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		form := url.Values{}
		form.Set("id", "t1_test")

		_ = r.ParseForm()
		if expect, actual := form, r.PostForm; !reflect.DeepEqual(expect, actual) {
			t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
		}

		fmt.Fprint(w, `{}`)
	})

	_, err := client.Comment.Unsave(ctx, "t3_test")
	if err == nil {
		t.Fatal("expected error, got nothing instead")
	}
	if expect, actual := `must provide comment id (starting with t1_); id provided: "t3_test"`, err.Error(); expect != actual {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}

	res, err := client.Comment.Unsave(ctx, "t1_test")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	if expect, actual := http.StatusOK, res.StatusCode; expect != actual {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}
}
