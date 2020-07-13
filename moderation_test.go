package reddit

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var expectedModActionsResult = &ModActions{
	ModActions: []ModAction{
		{
			ID:      "ModAction_b4e7979a-c4ad-11ea-8440-0ea1b7c2b8f9",
			Action:  "spamcomment",
			Created: &Timestamp{time.Date(2020, 7, 13, 2, 8, 14, 0, time.UTC)},

			Moderator:   "v_95",
			ModeratorID: "164ab8",

			TargetAuthor:    "testuser",
			TargetID:        "t1_fxw10aa",
			TargetPermalink: "/r/helloworldtestt/comments/hq6r3t/yo/fxw10aa/",
			TargetBody:      "hi",

			Subreddit:   "helloworldtestt",
			SubredditID: "2uquw1",
		},
		{
			ID:      "ModAction_a0408162-c4ad-11ea-8239-0e3b48262e8b",
			Action:  "sticky",
			Created: &Timestamp{time.Date(2020, 7, 13, 2, 7, 38, 0, time.UTC)},

			Moderator:   "v_95",
			ModeratorID: "164ab8",

			TargetAuthor:    "testuser",
			TargetID:        "t3_hq6r3t",
			TargetTitle:     "yo",
			TargetPermalink: "/r/helloworldtestt/comments/hq6r3t/yo/",

			Subreddit:   "helloworldtestt",
			SubredditID: "2uquw1",
		},
	},
	After:  "ModAction_a0408162-c4ad-11ea-8239-0e3b48262e8b",
	Before: "",
}

func TestModerationService_GetActions(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/moderation/actions.json")

	mux.HandleFunc("/r/testsubreddit/about/log", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	result, _, err := client.Moderation.GetActions(ctx, "testsubreddit")
	assert.NoError(t, err)
	assert.Equal(t, expectedModActionsResult, result)
}
