package geddit

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedFlairs = []Flair{
	{
		ID:   "b8a1c822-3feb-11e8-88e1-0e5f55d58ce0",
		Text: "Beginner",
		Type: "text",
		CSS:  "Beginner1",
	},
	{
		ID:   "b8ea0fce-3feb-11e8-af7a-0e263a127cf8",
		Text: "Beginner",
		Type: "text",
		CSS:  "Beginner2",
	},
}

var expectedFlairsV2 = []FlairV2{
	{
		ID:      "b8a1c822-3feb-11e8-88e1-0e5f55d58ce0",
		Text:    "Beginner",
		Type:    "text",
		CSS:     "Beginner1",
		ModOnly: false,
	},
	{
		ID:      "b8ea0fce-3feb-11e8-af7a-0e263a127cf8",
		Text:    "Moderator",
		Type:    "text",
		CSS:     "Moderator1",
		ModOnly: true,
	},
}

func TestFlairServiceOp_GetFlairs(t *testing.T) {
	setup()
	defer teardown()

	flairsBlob := readFileContents(t, "testdata/flairs.json")

	mux.HandleFunc("/r/subreddit/api/user_flair", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, flairsBlob)
	})

	flairs, _, err := client.Flair.GetFromSubreddit(ctx, "subreddit")
	assert.NoError(t, err)
	assert.Equal(t, expectedFlairs, flairs)
}

func TestFlairServiceOp_GetFlairsV2(t *testing.T) {
	setup()
	defer teardown()

	flairsV2Blob := readFileContents(t, "testdata/flairs-v2.json")

	mux.HandleFunc("/r/subreddit/api/user_flair_v2", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, flairsV2Blob)
	})

	flairs, _, err := client.Flair.GetFromSubredditV2(ctx, "subreddit")
	assert.NoError(t, err)
	assert.Equal(t, expectedFlairsV2, flairs)
}
