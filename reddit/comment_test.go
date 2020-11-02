package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedCommentSubmitOrEdit = &Comment{
	ID:        "test2",
	FullID:    "t1_test2",
	ParentID:  "t1_test",
	Permalink: "/r/subreddit/comments/test1/some_thread/test2/",

	Body:            "test comment",
	Author:          "reddit_username",
	AuthorID:        "t2_user1",
	AuthorFlairText: "Flair",
	AuthorFlairID:   "024b2b66-05ca-11e1-96f4-12313d096aae",

	SubredditName:         "subreddit",
	SubredditNamePrefixed: "r/subreddit",
	SubredditID:           "t5_test",

	Likes: Bool(true),

	Score:            1,
	Controversiality: 0,

	Created: &Timestamp{time.Date(2020, 4, 29, 0, 9, 47, 0, time.UTC)},
	Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

	PostID: "t3_link1",
}

func TestCommentService_Submit(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/comment/submit-or-edit.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/comment", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("return_rtjson", "true")
		form.Set("parent", "t1_test")
		form.Set("text", "test comment")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	comment, _, err := client.Comment.Submit(ctx, "t1_test", "test comment")
	require.NoError(t, err)
	require.Equal(t, expectedCommentSubmitOrEdit, comment)
}

func TestCommentService_Edit(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/comment/submit-or-edit.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/editusertext", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("return_rtjson", "true")
		form.Set("thing_id", "t1_test")
		form.Set("text", "test comment")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	comment, _, err := client.Comment.Edit(ctx, "t1_test", "test comment")
	require.NoError(t, err)
	require.Equal(t, expectedCommentSubmitOrEdit, comment)
}

func TestCommentService_Delete(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/del", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Comment.Delete(ctx, "t1_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCommentService_Save(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/save", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Comment.Save(ctx, "t1_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCommentService_Unsave(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/unsave", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Comment.Unsave(ctx, "t1_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCommentService_EnableReplies(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/sendreplies", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")
		form.Set("state", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Comment.EnableReplies(ctx, "t1_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCommentService_DisableReplies(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/sendreplies", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")
		form.Set("state", "false")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Comment.DisableReplies(ctx, "t1_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCommentService_Lock(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/lock", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Comment.Lock(ctx, "t1_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCommentService_Unlock(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/unlock", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Comment.Unlock(ctx, "t1_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCommentService_Upvote(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")
		form.Set("dir", "1")
		form.Set("rank", "10")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Comment.Upvote(ctx, "t1_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCommentService_Downvote(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")
		form.Set("dir", "-1")
		form.Set("rank", "10")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Comment.Downvote(ctx, "t1_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCommentService_RemoveVote(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")
		form.Set("dir", "0")
		form.Set("rank", "10")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Comment.RemoveVote(ctx, "t1_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCommentService_LoadMoreReplies(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/comment/more.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/morechildren", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("link_id", "t3_123")
		form.Set("children", "def,ghi,jkl")
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	_, err = client.Comment.LoadMoreReplies(ctx, nil)
	require.EqualError(t, err, "*Comment: cannot be nil")

	resp, err := client.Comment.LoadMoreReplies(ctx, &Comment{})
	require.Nil(t, resp)
	require.Nil(t, err)

	comment := &Comment{
		FullID: "t1_abc",
		PostID: "t3_123",
		Replies: Replies{
			More: &More{
				Children: []string{"def", "ghi", "jkl"},
			},
		},
	}

	_, err = client.Comment.LoadMoreReplies(ctx, comment)
	require.Nil(t, err)
	require.False(t, comment.HasMore())
	require.Len(t, comment.Replies.Comments, 2)
	require.Len(t, comment.Replies.Comments[0].Replies.Comments, 1)
}

func TestCommentService_Report(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/report", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("thing_id", "t1_test")
		form.Set("reason", "test reason")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Comment.Report(ctx, "t1_test", "test reason")
	require.NoError(t, err)
}
