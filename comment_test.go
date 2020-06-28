package geddit

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

func TestCommentServiceOp_Submit(t *testing.T) {
	setup()
	defer teardown()

	commentBlob := readFileContents(t, "testdata/comment-submit-edit.json")

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

		fmt.Fprint(w, commentBlob)
	})

	comment, _, err := client.Comment.Submit(ctx, "t1_test", "test comment")
	assert.NoError(t, err)
	assert.Equal(t, expectedCommentSubmitOrEdit, comment)
}

func TestCommentServiceOp_Edit(t *testing.T) {
	setup()
	defer teardown()

	commentBlob := readFileContents(t, "testdata/comment-submit-edit.json")

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

		fmt.Fprint(w, commentBlob)
	})

	_, _, err := client.Comment.Edit(ctx, "t3_test", "test comment")
	assert.EqualError(t, err, `must provide comment id (starting with t1_); id provided: "t3_test"`)

	comment, _, err := client.Comment.Edit(ctx, "t1_test", "test comment")
	assert.NoError(t, err)
	assert.Equal(t, expectedCommentSubmitOrEdit, comment)
}

func TestCommentServiceOp_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/del", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{}`)
	})

	_, err := client.Comment.Delete(ctx, "t3_test")
	assert.EqualError(t, err, `must provide comment id (starting with t1_); id provided: "t3_test"`)

	res, err := client.Comment.Delete(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestCommentServiceOp_Save(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/save", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{}`)
	})

	_, err := client.Comment.Save(ctx, "t3_test")
	assert.EqualError(t, err, `must provide comment id (starting with t1_); id provided: "t3_test"`)

	res, err := client.Comment.Save(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestCommentServiceOp_Unsave(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unsave", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{}`)
	})

	_, err := client.Comment.Unsave(ctx, "t3_test")
	assert.EqualError(t, err, `must provide comment id (starting with t1_); id provided: "t3_test"`)

	res, err := client.Comment.Unsave(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
