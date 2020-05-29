package geddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoteServiceOp_Up(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("dir", fmt.Sprint(upvote))
		form.Set("rank", "10")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{}`)
	})

	res, err := client.Vote.Up(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestVoteServiceOp_Down(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("dir", fmt.Sprint(downvote))
		form.Set("rank", "10")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{}`)
	})

	res, err := client.Vote.Down(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestVoteServiceOp_Remove(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("dir", fmt.Sprint(novote))
		form.Set("rank", "10")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{}`)
	})

	res, err := client.Vote.Remove(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
