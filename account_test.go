package geddit

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedKarma = []SubredditKarma{
	{Subreddit: "nba", PostKarma: 144, CommentKarma: 21999},
	{Subreddit: "redditdev", PostKarma: 19, CommentKarma: 4},
	{Subreddit: "test", PostKarma: 1, CommentKarma: 0},
	{Subreddit: "golang", PostKarma: 1, CommentKarma: 0},
}

func TestAccountServiceOp_Karma(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/account/karma.json")

	mux.HandleFunc("/api/v1/me/karma", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	karma, _, err := client.Account.Karma(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedKarma, karma)
}
