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
		form.Set("id", "t3_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.MarkNSFW(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_UnmarkNSFW(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unmarknsfw", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.UnmarkNSFW(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Spoiler(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/spoiler", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Spoiler(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Unspoiler(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unspoiler", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Unspoiler(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Sticky(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_subreddit_sticky", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "true")
		form.Set("num", "1")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Sticky(ctx, "t3_test", false)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Unsticky(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_subreddit_sticky", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "false")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Unsticky(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_PinToProfile(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_subreddit_sticky", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "true")
		form.Set("to_profile", "true")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.PinToProfile(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_UnpinFromProfile(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_subreddit_sticky", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "false")
		form.Set("to_profile", "true")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.UnpinFromProfile(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
