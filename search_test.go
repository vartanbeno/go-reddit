package reddit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var expectedSearchPosts = &Posts{
	Posts: []*Post{
		{
			ID:      "hybow9",
			FullID:  "t3_hybow9",
			Created: &Timestamp{time.Date(2020, 7, 26, 18, 14, 24, 0, time.UTC)},
			Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

			Permalink: "https://www.reddit.com/r/WatchPeopleDieInside/comments/hybow9/pregnancy_test/",
			URL:       "https://v.redd.it/ra4qnt8bt8d51",

			Title: "Pregnancy test",

			Score:            103829,
			UpvoteRatio:      0.88,
			NumberOfComments: 3748,

			SubredditID:           "t5_3h4zq",
			SubredditName:         "WatchPeopleDieInside",
			SubredditNamePrefixed: "r/WatchPeopleDieInside",

			AuthorID:   "t2_3p32m02",
			AuthorName: "chocolat_ice_cream",
		},
		{
			ID:      "hmwhd7",
			FullID:  "t3_hmwhd7",
			Created: &Timestamp{time.Date(2020, 7, 7, 15, 19, 42, 0, time.UTC)},
			Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

			Permalink: "https://www.reddit.com/r/worldnews/comments/hmwhd7/brazilian_president_jair_bolsonaro_tests_positive/",
			URL:       "https://www.theguardian.com/world/2020/jul/07/jair-bolsonaro-coronavirus-positive-test-brazil-president",

			Title: "Brazilian president Jair Bolsonaro tests positive for coronavirus",

			Score:            149238,
			UpvoteRatio:      0.94,
			NumberOfComments: 7415,

			SubredditID:           "t5_2qh13",
			SubredditName:         "worldnews",
			SubredditNamePrefixed: "r/worldnews",

			AuthorID:   "t2_wgrkg",
			AuthorName: "Jeremy_Martin",
		},
	},
	After: "t3_hmwhd7",
}

var expectedSearchUsers = &Users{
	Users: []*User{
		{
			ID:      "abc",
			Name:    "user1",
			Created: &Timestamp{time.Date(2019, 8, 14, 23, 38, 42, 0, time.UTC)},

			PostKarma:    5730,
			CommentKarma: 11740,

			HasVerifiedEmail: true,
		},
		{
			ID:      "def",
			Name:    "user2",
			Created: &Timestamp{time.Date(2020, 5, 7, 3, 16, 46, 0, time.UTC)},

			PostKarma:    2485,
			CommentKarma: 127,
		},
	},
}

var expectedSearchSubreddits = &Subreddits{
	Subreddits: []*Subreddit{
		{
			ID:      "2qh23",
			FullID:  "t5_2qh23",
			Created: &Timestamp{time.Date(2008, 1, 25, 5, 11, 28, 0, time.UTC)},

			URL:          "/r/test/",
			Name:         "test",
			NamePrefixed: "r/test",
			Title:        "Testing",
			Type:         "public",

			Subscribers: 8174,
		},
		{
			ID:      "333yu",
			FullID:  "t5_333yu",
			Created: &Timestamp{time.Date(2014, 8, 18, 23, 29, 47, 0, time.UTC)},

			URL:          "/r/trollingforababy/",
			Name:         "trollingforababy",
			NamePrefixed: "r/trollingforababy",
			Title:        "Crushing it with reddit karma",
			Description:  "This is a group for laughing at and mocking the awkward, ridiculous, and sometimes painful things we endure while trying for a baby. Trollingforababy is for people who are trying to conceive, and are not currently pregnant. \n\nPlease look at our complete list of rules before participating.",
			Type:         "public",

			Subscribers: 10244,
		},
	},
	After: "t5_333yu",
}

func TestSearchService_Posts(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/search/posts.json")
	assert.NoError(t, err)

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("type", "link")
		form.Set("q", "test")
		form.Set("after", "t3_testpost")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	posts, _, err := client.Search.Posts(ctx, "test", nil, SetAfter("t3_testpost"))
	assert.NoError(t, err)
	assert.Equal(t, expectedSearchPosts, posts)
}

func TestSearchService_Posts_InSubreddit(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/search/posts.json")
	assert.NoError(t, err)

	mux.HandleFunc("/r/test/search", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("type", "link")
		form.Set("q", "test")
		form.Set("restrict_sr", "true")
		form.Set("after", "t3_testpost")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	posts, _, err := client.Search.Posts(ctx, "test", []string{"test"}, SetAfter("t3_testpost"))
	assert.NoError(t, err)
	assert.Equal(t, expectedSearchPosts, posts)
}

func TestSearchService_Posts_InSubreddits(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/search/posts.json")
	assert.NoError(t, err)

	mux.HandleFunc("/r/test+golang+nba/search", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("type", "link")
		form.Set("q", "test")
		form.Set("restrict_sr", "true")
		form.Set("after", "t3_testpost")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	posts, _, err := client.Search.Posts(ctx, "test", []string{"test", "golang", "nba"}, SetAfter("t3_testpost"))
	assert.NoError(t, err)
	assert.Equal(t, expectedSearchPosts, posts)
}

func TestSearchService_Subreddits(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/search/subreddits.json")
	assert.NoError(t, err)

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("type", "sr")
		form.Set("q", "test")
		form.Set("before", "t5_testsr")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	subreddits, _, err := client.Search.Subreddits(ctx, "test", SetBefore("t5_testsr"))
	assert.NoError(t, err)
	assert.Equal(t, expectedSearchSubreddits, subreddits)
}

func TestSearchService_Users(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/search/users.json")
	assert.NoError(t, err)

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("type", "user")
		form.Set("q", "test")
		form.Set("limit", "2")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	users, _, err := client.Search.Users(ctx, "test", SetLimit(2))
	assert.NoError(t, err)
	assert.Equal(t, expectedSearchUsers, users)
}

func printJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
