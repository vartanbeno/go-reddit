package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var expectedCommentSubmitOrEdit = &Comment{
	ID:        "test2",
	FullID:    "t1_test2",
	ParentID:  "t1_test",
	Permalink: "https://www.reddit.com/r/subreddit/comments/test1/some_thread/test2/",

	Body:            "test comment",
	Author:          "reddit_username",
	AuthorID:        "t2_user1",
	AuthorFlairText: "Flair",
	AuthorFlairID:   "024b2b66-05ca-11e1-96f4-12313d096aae",

	Subreddit:             "subreddit",
	SubredditNamePrefixed: "r/subreddit",
	SubredditID:           "t5_test",

	Likes: Bool(true),

	Score:            1,
	Controversiality: 0,

	Created: &Timestamp{time.Date(2020, 4, 29, 0, 9, 47, 0, time.UTC)},
	// todo: this should just be nil
	Edited: &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

	PostID: "t3_link1",
}

func TestCommentService_Submit(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/comment-submit-edit.json")

	mux.HandleFunc("/api/comment", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("return_rtjson", "true")
		form.Set("parent", "t1_test")
		form.Set("text", "test comment")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	comment, _, err := client.Comment.Submit(ctx, "t1_test", "test comment")
	assert.NoError(t, err)
	assert.Equal(t, expectedCommentSubmitOrEdit, comment)
}

func TestCommentService_Edit(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/comment-submit-edit.json")

	mux.HandleFunc("/api/editusertext", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("return_rtjson", "true")
		form.Set("thing_id", "t1_test")
		form.Set("text", "test comment")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	comment, _, err := client.Comment.Edit(ctx, "t1_test", "test comment")
	assert.NoError(t, err)
	assert.Equal(t, expectedCommentSubmitOrEdit, comment)
}
