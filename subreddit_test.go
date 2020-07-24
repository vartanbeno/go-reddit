package reddit

import (
	"fmt"
	"net/http"
	"net/url"
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
	Subreddits: []*Subreddit{
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
			Favorite:    false,
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
			Favorite:    true,
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
			Favorite:    false,
		},
	},
}

var expectedSticky = &postAndComments{
	Post: &Post{
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

var expectSubredditInfos = []*SubredditInfo{
	{Name: "golang", Subscribers: 119_722, ActiveUsers: 531},
	{Name: "golang_infosec", Subscribers: 1_776, ActiveUsers: 0},
	{Name: "GolangJobOfferings", Subscribers: 863, ActiveUsers: 1},
	{Name: "golang2", Subscribers: 626, ActiveUsers: 0},
	{Name: "GolangUnofficial", Subscribers: 239, ActiveUsers: 4},
	{Name: "golanguage", Subscribers: 247, ActiveUsers: 4},
	{Name: "golang_jobs", Subscribers: 16, ActiveUsers: 4},
}

var expectSubredditNames = []string{
	"golang",
	"golang_infosec",
	"GolangJobOfferings",
	"golanguage",
	"golang2",
	"GolangUnofficial",
	"golang_jobs",
}

var expectedModerators = []Moderator{
	{ID: "t2_test1", Name: "testuser1", Permissions: []string{"all"}},
	{ID: "t2_test2", Name: "testuser2", Permissions: []string{"all"}},
}

var expectedRandomSubreddit = &Subreddit{
	FullID:  "t5_2wi4l",
	Created: &Timestamp{time.Date(2013, 3, 1, 4, 4, 18, 0, time.UTC)},

	URL:          "/r/GalaxyS8/",
	Name:         "GalaxyS8",
	NamePrefixed: "r/GalaxyS8",
	Title:        "Samsung Galaxy S8",
	Description:  "The only place for news, discussion, photos, and everything else Samsung Galaxy S8.",
	Type:         "public",

	Subscribers: 52357,
}

func TestSubredditService_Get(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/subreddit/about.json")
	assert.NoError(t, err)

	mux.HandleFunc("/r/golang/about", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	_, _, err = client.Subreddit.Get(ctx, "")
	assert.EqualError(t, err, "name: must not be empty")

	subreddit, _, err := client.Subreddit.Get(ctx, "golang")
	assert.NoError(t, err)
	assert.Equal(t, expectedSubreddit, subreddit)
}

func TestSubredditService_GetPopular(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/subreddit/list.json")
	assert.NoError(t, err)

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

	blob, err := readFileContents("testdata/subreddit/list.json")
	assert.NoError(t, err)

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

	blob, err := readFileContents("testdata/subreddit/list.json")
	assert.NoError(t, err)

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

	blob, err := readFileContents("testdata/subreddit/list.json")
	assert.NoError(t, err)

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

	blob, err := readFileContents("testdata/subreddit/list.json")
	assert.NoError(t, err)

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

	blob, err := readFileContents("testdata/subreddit/list.json")
	assert.NoError(t, err)

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

	blob, err := readFileContents("testdata/subreddit/list.json")
	assert.NoError(t, err)

	mux.HandleFunc("/subreddits/mine/moderator", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, _, err := client.Subreddit.GetModerated(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedSubreddits, subreddits)
}

func TestSubredditService_GetSticky1(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/post/post.json")
	assert.NoError(t, err)

	mux.HandleFunc("/r/test/about/sticky", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "1", r.Form.Get("num"))

		fmt.Fprint(w, blob)
	})

	post, comments, _, err := client.Subreddit.GetSticky1(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, expectedPost2, post)
	assert.Equal(t, expectedComments, comments)
}

func TestSubredditService_GetSticky2(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/post/post.json")
	assert.NoError(t, err)

	mux.HandleFunc("/r/test/about/sticky", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "2", r.Form.Get("num"))

		fmt.Fprint(w, blob)
	})

	post, comments, _, err := client.Subreddit.GetSticky2(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, expectedPost2, post)
	assert.Equal(t, expectedComments, comments)
}

func TestSubredditService_Subscribe(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("action", "sub")
		form.Set("sr_name", "test,golang,nba")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)
	})

	_, err := client.Subreddit.Subscribe(ctx, "test", "golang", "nba")
	assert.NoError(t, err)
}

func TestSubredditService_SubscribeByID(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("action", "sub")
		form.Set("sr", "t5_test1,t5_test2,t5_test3")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)
	})

	_, err := client.Subreddit.SubscribeByID(ctx, "t5_test1", "t5_test2", "t5_test3")
	assert.NoError(t, err)
}

func TestSubredditService_Unsubscribe(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("action", "unsub")
		form.Set("sr_name", "test,golang,nba")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)
	})

	_, err := client.Subreddit.Unsubscribe(ctx, "test", "golang", "nba")
	assert.NoError(t, err)
}

func TestSubredditService_UnsubscribeByID(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("action", "unsub")
		form.Set("sr", "t5_test1,t5_test2,t5_test3")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)
	})

	_, err := client.Subreddit.UnsubscribeByID(ctx, "t5_test1", "t5_test2", "t5_test3")
	assert.NoError(t, err)
}

func TestSubredditService_Search(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/subreddit/search.json")
	assert.NoError(t, err)

	mux.HandleFunc("/api/search_subreddits", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("query", "golang")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	subreddits, _, err := client.Subreddit.Search(ctx, "golang")
	assert.NoError(t, err)
	assert.Equal(t, expectSubredditInfos, subreddits)
}

func TestSubredditService_SearchNames(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/subreddit/search-names.json")
	assert.NoError(t, err)

	mux.HandleFunc("/api/search_reddit_names", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("query", "golang")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	names, _, err := client.Subreddit.SearchNames(ctx, "golang")
	assert.NoError(t, err)
	assert.Equal(t, expectSubredditNames, names)
}

func TestSubredditService_Moderators(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/subreddit/moderators.json")
	assert.NoError(t, err)

	mux.HandleFunc("/r/test/about/moderators", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	moderators, _, err := client.Subreddit.Moderators(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, expectedModerators, moderators)
}

func TestSubredditService_Random(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/subreddit/random.json")
	assert.NoError(t, err)

	mux.HandleFunc("/r/random", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "true", r.Form.Get("sr_detail"))
		assert.Equal(t, "1", r.Form.Get("limit"))

		fmt.Fprint(w, blob)
	})

	subreddit, _, err := client.Subreddit.Random(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedRandomSubreddit, subreddit)
}

func TestSubredditService_RandomNSFW(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/subreddit/random.json")
	assert.NoError(t, err)

	mux.HandleFunc("/r/randnsfw", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "true", r.Form.Get("sr_detail"))
		assert.Equal(t, "1", r.Form.Get("limit"))

		fmt.Fprint(w, blob)
	})

	subreddit, _, err := client.Subreddit.RandomNSFW(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedRandomSubreddit, subreddit)
}
