package geddit

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
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
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, flairsBlob)
	})

	flairs, _, err := client.Flair.GetFromSubreddit(ctx, "subreddit")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	if expect, actual := expectedFlairs, flairs; !reflect.DeepEqual(expect, actual) {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}
}

func TestFlairServiceOp_GetFlairsV2(t *testing.T) {
	setup()
	defer teardown()

	flairsV2Blob := readFileContents(t, "testdata/flairs-v2.json")

	mux.HandleFunc("/r/subreddit/api/user_flair_v2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, flairsV2Blob)
	})

	flairs, _, err := client.Flair.GetFromSubredditV2(ctx, "subreddit")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	if expect, actual := expectedFlairsV2, flairs; !reflect.DeepEqual(expect, actual) {
		t.Fatalf("got unexpected value\nexpect: %s\nactual: %s", Stringify(expect), Stringify(actual))
	}
}
