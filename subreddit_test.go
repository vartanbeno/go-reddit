package geddit

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var expectedSubreddit = &Subreddit{
	ID:      "2rc7j",
	FullID:  "t5_2rc7j",
	Created: &Timestamp{time.Date(2009, 11, 11, 0, 54, 28, 0, time.UTC)},

	URL:          "/r/golang/",
	Name:         "golang",
	NamePrefixed: "r/golang",
	Title:        "The Go Programming Language",
	Description:  "Ask questions and post articles about the Go programming language and related tools, events etc.",
	Type:         "public",

	Subscribers:     116532,
	ActiveUserCount: Int(386),
	NSFW:            false,
	UserIsMod:       false,
}

var expectedSubreddits = &Subreddits{
	After:  "t5_2qh0u",
	Before: "",
	Subreddits: []Subreddit{
		{
			ID:      "2qs0k",
			FullID:  "t5_2qs0k",
			Created: &Timestamp{time.Date(2009, 1, 25, 2, 25, 57, 0, time.UTC)},

			URL:          "/r/Home/",
			Name:         "Home",
			NamePrefixed: "r/Home",
			Title:        "Home",
			Type:         "public",

			Subscribers: 15336,
			NSFW:        false,
			UserIsMod:   false,
		},
		{
			ID:      "2qh1i",
			FullID:  "t5_2qh1i",
			Created: &Timestamp{time.Date(2008, 1, 25, 3, 52, 15, 0, time.UTC)},

			URL:          "/r/AskReddit/",
			Name:         "AskReddit",
			NamePrefixed: "r/AskReddit",
			Title:        "Ask Reddit...",
			Description:  "r/AskReddit is the place to ask and answer thought-provoking questions.",
			Type:         "public",

			Subscribers: 28449174,
			NSFW:        false,
			UserIsMod:   false,
		},
		{
			ID:      "2qh0u",
			FullID:  "t5_2qh0u",
			Created: &Timestamp{time.Date(2008, 1, 25, 0, 31, 9, 0, time.UTC)},

			URL:          "/r/pics/",
			Name:         "pics",
			NamePrefixed: "r/pics",
			Title:        "Reddit Pics",
			Description:  "A place for pictures and photographs.",
			Type:         "public",

			Subscribers: 24987753,
			NSFW:        false,
			UserIsMod:   false,
		},
	},
}

var expectedSticky = &PostAndComments{
	Post: Post{
		ID:      "hcl9gq",
		FullID:  "t3_hcl9gq",
		Created: &Timestamp{time.Date(2020, 6, 20, 12, 8, 57, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "https://www.reddit.com/r/nba/comments/hcl9gq/daily_discussion_thread_freetalk_and_other/",
		URL:       "https://www.reddit.com/r/nba/comments/hcl9gq/daily_discussion_thread_freetalk_and_other/",

		Title: "Daily Discussion Thread | Free-Talk and Other Updates - June 20, 2020",
		Body:  "Talk about whatever is on your mind, basketball related or not.\n\n# Useful Links \u0026amp; Other Resources\n\n[List of All #NBATogether Live Classic Games Streamed to Date](https://www.youtube.com/results?search_query=%23NBATogetherLive)\n\n[r/nba Discord Server](https://www.discord.gg/nba)\n\n[r/nba Twitter](https://twitter.com/nba_reddit)\n\n[Read Our Community's Rules and Guidelines](https://www.reddit.com/r/nba/wiki/rules)",

		Score:            16,
		UpvoteRatio:      0.82,
		NumberOfComments: 25,

		SubredditID:           "t5_2qo4s",
		SubredditName:         "nba",
		SubredditNamePrefixed: "r/nba",

		AuthorID:   "t2_6l4z3",
		AuthorName: "AutoModerator",

		IsSelfPost: true,
		Stickied:   true,
	},
}

func TestSubredditService_GetByName(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/subreddit/about.json")

	mux.HandleFunc("/r/golang/about", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddit, _, err := client.Subreddit.GetByName(ctx, "golang")
	assert.NoError(t, err)
	assert.Equal(t, expectedSubreddit, subreddit)
}

func TestSubredditService_GetPopular(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/subreddit/list.json")

	mux.HandleFunc("/subreddits/popular", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, _, err := client.Subreddit.GetPopular(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedSubreddits, subreddits)
}

func TestSubredditService_GetNew(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/subreddit/list.json")

	mux.HandleFunc("/subreddits/new", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, _, err := client.Subreddit.GetNew(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedSubreddits, subreddits)
}

func TestSubredditService_GetGold(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/subreddit/list.json")

	mux.HandleFunc("/subreddits/gold", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, _, err := client.Subreddit.GetGold(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedSubreddits, subreddits)
}

func TestSubredditService_GetDefault(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/subreddit/list.json")

	mux.HandleFunc("/subreddits/default", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, _, err := client.Subreddit.GetDefault(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedSubreddits, subreddits)
}

func TestSubredditService_GetSubscribed(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/subreddit/list.json")

	mux.HandleFunc("/subreddits/mine/subscriber", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, _, err := client.Subreddit.GetSubscribed(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedSubreddits, subreddits)
}

func TestSubredditService_GetApproved(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/subreddit/list.json")

	mux.HandleFunc("/subreddits/mine/contributor", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, _, err := client.Subreddit.GetApproved(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedSubreddits, subreddits)
}

func TestSubredditService_GetModerated(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/subreddit/list.json")

	mux.HandleFunc("/subreddits/mine/moderator", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, _, err := client.Subreddit.GetModerated(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedSubreddits, subreddits)
}

// todo: WIP
func TestSubredditService_GetSticky1(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/subreddit/sticky.json")

	mux.HandleFunc("/r/nba/about/sticky", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "1", r.Form.Get("num"))

		fmt.Fprint(w, blob)
	})

	sticky, _, err := client.Subreddit.GetSticky1(ctx, "nba")
	assert.NoError(t, err)
	assert.Equal(t, expectedSticky.Post, sticky.Post)
	// b, _ := json.MarshalIndent(sticky.Comments, "", "  ")
	// fmt.Println(string(b))
}
