package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostAndCommentService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/del", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.PostAndComment.Delete(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostAndCommentService_Save(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/save", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.PostAndComment.Save(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostAndCommentService_Unsave(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unsave", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.PostAndComment.Unsave(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostAndCommentService_EnableReplies(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/sendreplies", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")
		form.Set("state", "true")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.PostAndComment.EnableReplies(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostAndCommentService_DisableReplies(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/sendreplies", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")
		form.Set("state", "false")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.PostAndComment.DisableReplies(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostAndCommentService_Lock(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/lock", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.PostAndComment.Lock(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostAndCommentService_Unlock(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unlock", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.PostAndComment.Unlock(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostAndCommentService_Upvote(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")
		form.Set("dir", fmt.Sprint(upvote))
		form.Set("rank", "10")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.PostAndComment.Upvote(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostAndCommentService_Downvote(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")
		form.Set("dir", fmt.Sprint(downvote))
		form.Set("rank", "10")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.PostAndComment.Downvote(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostAndCommentService_RemoveVote(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")
		form.Set("dir", fmt.Sprint(novote))
		form.Set("rank", "10")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.PostAndComment.RemoveVote(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
