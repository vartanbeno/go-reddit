package reddit

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostService_Hide(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/hide", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "1,2,3")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	_, err := client.Post.Hide(ctx)
	assert.EqualError(t, err, "must provide at least 1 id")

	res, err := client.Post.Hide(ctx, "1", "2", "3")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Unhide(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unhide", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "1,2,3")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	_, err := client.Post.Unhide(ctx)
	assert.EqualError(t, err, "must provide at least 1 id")

	res, err := client.Post.Unhide(ctx, "1", "2", "3")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_MarkNSFW(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/marknsfw", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.MarkNSFW(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_UnmarkNSFW(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unmarknsfw", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.UnmarkNSFW(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Spoiler(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/spoiler", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Spoiler(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Unspoiler(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unspoiler", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t1_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Unspoiler(ctx, "t1_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
