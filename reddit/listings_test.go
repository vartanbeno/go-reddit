package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedListingPosts = []*Post{
	{
		ID:      "i2gvg4",
		FullID:  "t3_i2gvg4",
		Created: &Timestamp{time.Date(2020, 8, 2, 18, 23, 8, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/test/comments/i2gvg4/this_is_a_title/",
		URL:       "https://www.reddit.com/r/test/comments/i2gvg4/this_is_a_title/",

		Title: "This is a title",
		Body:  "This is some text",

		Likes:            Bool(true),
		Score:            1,
		UpvoteRatio:      1,
		NumberOfComments: 1,

		SubredditName:         "test",
		SubredditNamePrefixed: "r/test",
		SubredditID:           "t5_2qh23",
		SubredditSubscribers:  8202,

		Author:   "v_95",
		AuthorID: "t2_164ab8",

		IsSelfPost: true,
	},
}

var expectedListingComments = []*Comment{
	{
		ID:      "g05v931",
		FullID:  "t1_g05v931",
		Created: &Timestamp{time.Date(2020, 8, 3, 1, 15, 40, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		ParentID:  "t3_i2gvg4",
		Permalink: "/r/test/comments/i2gvg4/this_is_a_title/g05v931/",

		Body:     "Test comment",
		Author:   "v_95",
		AuthorID: "t2_164ab8",

		SubredditName:         "test",
		SubredditNamePrefixed: "r/test",
		SubredditID:           "t5_2qh23",

		Likes: Bool(true),

		Score:            1,
		Controversiality: 0,

		PostID: "t3_i2gvg4",

		IsSubmitter: true,
	},
}

var expectedListingSubreddits = []*Subreddit{
	{
		ID:      "2qh23",
		FullID:  "t5_2qh23",
		Created: &Timestamp{time.Date(2008, 1, 25, 5, 11, 28, 0, time.UTC)},

		URL:          "/r/test/",
		Name:         "test",
		NamePrefixed: "r/test",
		Title:        "Testing",
		Type:         "public",

		Subscribers: 8202,
		Subscribed:  true,
	},
}

var expectedListingPosts2 = []*Post{
	{
		ID:      "i2gvg4",
		FullID:  "t3_i2gvg4",
		Created: &Timestamp{time.Date(2020, 8, 2, 18, 23, 8, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/test/comments/i2gvg4/this_is_a_title/",
		URL:       "https://www.reddit.com/r/test/comments/i2gvg4/this_is_a_title/",

		Title: "This is a title",
		Body:  "This is some text",

		Likes:            Bool(true),
		Score:            1,
		UpvoteRatio:      1,
		NumberOfComments: 1,

		SubredditName:         "test",
		SubredditNamePrefixed: "r/test",
		SubredditID:           "t5_2qh23",
		SubredditSubscribers:  8201,

		Author:   "v_95",
		AuthorID: "t2_164ab8",

		IsSelfPost: true,
	},
	{
		ID:      "i2gvs1",
		FullID:  "t3_i2gvs1",
		Created: &Timestamp{time.Date(2020, 8, 2, 18, 23, 37, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/test/comments/i2gvs1/this_is_a_title/",
		URL:       "http://example.com",

		Title: "This is a title",

		Likes:            Bool(true),
		Score:            1,
		UpvoteRatio:      1,
		NumberOfComments: 0,

		SubredditName:         "test",
		SubredditNamePrefixed: "r/test",
		SubredditID:           "t5_2qh23",
		SubredditSubscribers:  8201,

		Author:   "v_95",
		AuthorID: "t2_164ab8",
	},
}

func TestListingsService_Get(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/listings/posts-comments-subreddits.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("id", "t5_2qh23,t3_i2gvg4,t1_g05v931")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	posts, comments, subreddits, _, err := client.Listings.Get(ctx, "t5_2qh23", "t3_i2gvg4", "t1_g05v931")
	require.NoError(t, err)
	require.Equal(t, expectedListingPosts, posts)
	require.Equal(t, expectedListingComments, comments)
	require.Equal(t, expectedListingSubreddits, subreddits)
}

func TestListingsService_GetPosts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/listings/posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/by_id/t3_i2gvg4,t3_i2gwgz", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, _, err := client.Listings.GetPosts(ctx, "t3_i2gvg4", "t3_i2gwgz")
	require.NoError(t, err)
	require.Equal(t, expectedListingPosts2, posts)
}
