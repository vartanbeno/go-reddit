package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedPostAndComments = &PostAndComments{
	Post: &Post{
		ID:      "testpost",
		FullID:  "t3_testpost",
		Created: &Timestamp{time.Date(2020, 7, 18, 10, 26, 7, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/test/comments/testpost/test/",
		URL:       "https://www.reddit.com/r/test/comments/testpost/test/",

		Title: "Test",
		Body:  "Hello",

		Score:            1,
		UpvoteRatio:      1,
		NumberOfComments: 2,

		SubredditName:         "test",
		SubredditNamePrefixed: "r/test",
		SubredditID:           "t5_2qh23",
		SubredditSubscribers:  8077,

		Author:   "testuser",
		AuthorID: "t2_testuser",

		IsSelfPost: true,
	},
	Comments: []*Comment{
		{
			ID:      "testc1",
			FullID:  "t1_testc1",
			Created: &Timestamp{time.Date(2020, 7, 18, 10, 31, 59, 0, time.UTC)},
			Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

			ParentID:  "t3_testpost",
			Permalink: "/r/test/comments/testpost/test/testc1/",

			Body:     "Hi",
			Author:   "testuser",
			AuthorID: "t2_testuser",

			SubredditName:         "test",
			SubredditNamePrefixed: "r/test",
			SubredditID:           "t5_2qh23",

			Score:            1,
			Controversiality: 0,

			PostID: "t3_testpost",

			IsSubmitter: true,
			CanGild:     true,

			Replies: Replies{
				Comments: []*Comment{
					{
						ID:      "testc2",
						FullID:  "t1_testc2",
						Created: &Timestamp{time.Date(2020, 7, 18, 10, 32, 28, 0, time.UTC)},
						Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

						ParentID:  "t1_testc1",
						Permalink: "/r/test/comments/testpost/test/testc2/",

						Body:     "Hello",
						Author:   "testuser",
						AuthorID: "t2_testuser",

						SubredditName:         "test",
						SubredditNamePrefixed: "r/test",
						SubredditID:           "t5_2qh23",

						Score:            1,
						Controversiality: 0,

						PostID: "t3_testpost",

						IsSubmitter: true,
						CanGild:     true,
					},
				},
			},
		},
	},
}

var expectedSubmittedPost = &Submitted{
	ID:     "hw6l6a",
	FullID: "t3_hw6l6a",
	URL:    "https://www.reddit.com/r/test/comments/hw6l6a/test_title/",
}

var expectedEditedPost = &Post{
	ID:      "hw6l6a",
	FullID:  "t3_hw6l6a",
	Created: &Timestamp{time.Date(2020, 7, 23, 1, 24, 55, 0, time.UTC)},
	Edited:  &Timestamp{time.Date(2020, 7, 23, 1, 42, 44, 0, time.UTC)},

	Permalink: "/r/test/comments/hw6l6a/test_title/",
	URL:       "https://www.reddit.com/r/test/comments/hw6l6a/test_title/",

	Title: "Test Title",
	Body:  "this is edited",

	Likes: Bool(true),

	Score:            1,
	UpvoteRatio:      1,
	NumberOfComments: 0,

	SubredditName:         "test",
	SubredditNamePrefixed: "r/test",
	SubredditID:           "t5_2qh23",
	SubredditSubscribers:  8128,

	Author:   "v_95",
	AuthorID: "t2_164ab8",

	Spoiler:    true,
	IsSelfPost: true,
}

var expectedPost2 = &Post{
	ID:      "i2gvs1",
	FullID:  "t3_i2gvs1",
	Created: &Timestamp{time.Date(2020, 8, 2, 18, 23, 37, 0, time.UTC)},
	Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

	Permalink: "/r/test/comments/i2gvs1/this_is_a_title/",
	URL:       "http://example.com",

	Title: "This is a title",

	Likes: Bool(true),

	Score:            1,
	UpvoteRatio:      1,
	NumberOfComments: 0,

	SubredditName:         "test",
	SubredditNamePrefixed: "r/test",
	SubredditID:           "t5_2qh23",
	SubredditSubscribers:  8278,

	Author:   "v_95",
	AuthorID: "t2_164ab8",
}

var expectedPostDuplicates = []*Post{
	{
		ID:      "8kbs85",
		FullID:  "t3_8kbs85",
		Created: &Timestamp{time.Date(2018, 5, 18, 9, 10, 18, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/test/comments/8kbs85/test/",
		URL:       "http://example.com",

		Title: "test",

		Likes: nil,

		Score:            1,
		UpvoteRatio:      0.66,
		NumberOfComments: 1,

		SubredditName:         "test",
		SubredditNamePrefixed: "r/test",
		SubredditID:           "t5_2qh23",
		SubredditSubscribers:  8278,

		Author:   "GarlicoinAccount",
		AuthorID: "t2_d2v1r90",
	},
	{
		ID:      "le1tc",
		FullID:  "t3_le1tc",
		Created: &Timestamp{time.Date(2011, 10, 16, 13, 26, 40, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/test/comments/le1tc/test_to_see_if_this_fixes_the_problem_of_my_likes/",
		URL:       "http://www.example.com",

		Title: "Test to see if this fixes the problem of my \"likes\" from the last 7 months vanishing.",

		Likes: nil,

		Score:            2,
		UpvoteRatio:      1,
		NumberOfComments: 1,

		SubredditName:         "test",
		SubredditNamePrefixed: "r/test",
		SubredditID:           "t5_2qh23",
		SubredditSubscribers:  8278,

		Author:   "prog101",
		AuthorID: "t2_8dyo",
	},
}

func TestPostService_Get(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/post/post.json")
	require.NoError(t, err)

	mux.HandleFunc("/comments/abc123", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	postAndComments, _, err := client.Post.Get(ctx, "abc123")
	require.NoError(t, err)
	require.Equal(t, expectedPostAndComments, postAndComments)
}

func TestPostService_Duplicates(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/post/duplicates.json")
	require.NoError(t, err)

	mux.HandleFunc("/duplicates/abc123", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "2")
		form.Set("sr", "test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	post, postDuplicates, resp, err := client.Post.Duplicates(ctx, "abc123", &ListDuplicatePostOptions{
		ListOptions: ListOptions{
			Limit: 2,
		},
		Subreddit: "test",
	})
	require.NoError(t, err)
	require.Equal(t, expectedPost2, post)
	require.Equal(t, expectedPostDuplicates, postDuplicates)
	require.Equal(t, "t3_le1tc", resp.After)
}

func TestPostService_SubmitText(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/post/submit.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/submit", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("kind", "self")
		form.Set("sr", "test")
		form.Set("title", "Test Title")
		form.Set("text", "Test Text")
		form.Set("spoiler", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	submittedPost, _, err := client.Post.SubmitText(ctx, SubmitTextRequest{
		Subreddit: "test",
		Title:     "Test Title",
		Text:      "Test Text",
		Spoiler:   true,
	})
	require.NoError(t, err)
	require.Equal(t, expectedSubmittedPost, submittedPost)
}

func TestPostService_SubmitLink(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/post/submit.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/submit", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("kind", "link")
		form.Set("sr", "test")
		form.Set("title", "Test Title")
		form.Set("url", "https://www.example.com")
		form.Set("sendreplies", "false")
		form.Set("resubmit", "true")
		form.Set("nsfw", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	submittedPost, _, err := client.Post.SubmitLink(ctx, SubmitLinkRequest{
		Subreddit:   "test",
		Title:       "Test Title",
		URL:         "https://www.example.com",
		SendReplies: Bool(false),
		Resubmit:    true,
		NSFW:        true,
	})
	require.NoError(t, err)
	require.Equal(t, expectedSubmittedPost, submittedPost)
}

func TestPostService_Edit(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/post/edit.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/editusertext", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("return_rtjson", "true")
		form.Set("thing_id", "t3_test")
		form.Set("text", "test edit")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	editedPost, _, err := client.Post.Edit(ctx, "t3_test", "test edit")
	require.NoError(t, err)
	require.Equal(t, expectedEditedPost, editedPost)
}

func TestPostService_Hide(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/hide", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "1,2,3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Post.Hide(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	resp, err := client.Post.Hide(ctx, "1", "2", "3")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_Unhide(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/unhide", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "1,2,3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Post.Unhide(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	resp, err := client.Post.Unhide(ctx, "1", "2", "3")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_MarkNSFW(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/marknsfw", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.MarkNSFW(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_UnmarkNSFW(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/unmarknsfw", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.UnmarkNSFW(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_Spoiler(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/spoiler", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.Spoiler(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_Unspoiler(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/unspoiler", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.Unspoiler(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_Sticky(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_subreddit_sticky", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "true")
		form.Set("num", "1")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.Sticky(ctx, "t3_test", false)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_Unsticky(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_subreddit_sticky", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "false")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.Unsticky(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_PinToProfile(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_subreddit_sticky", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "true")
		form.Set("to_profile", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.PinToProfile(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_UnpinFromProfile(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_subreddit_sticky", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "false")
		form.Set("to_profile", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.UnpinFromProfile(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_SetSuggestedSortBest(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "confidence")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.SetSuggestedSortBest(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_SetSuggestedSortTop(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "top")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.SetSuggestedSortTop(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_SetSuggestedSortNew(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "new")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.SetSuggestedSortNew(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_SetSuggestedSortControversial(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "controversial")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.SetSuggestedSortControversial(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_SetSuggestedSortOld(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "old")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.SetSuggestedSortOld(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_SetSuggestedSortRandom(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "random")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.SetSuggestedSortRandom(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_SetSuggestedSortAMA(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "qa")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.SetSuggestedSortAMA(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_SetSuggestedSortLive(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "live")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.SetSuggestedSortLive(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_ClearSuggestedSort(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.ClearSuggestedSort(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_EnableContestMode(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_contest_mode", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.EnableContestMode(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_DisableContestMode(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/set_contest_mode", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "false")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.DisableContestMode(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_LoadMoreReplies(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/post/more.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/morechildren", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("link_id", "t3_123")
		form.Set("children", "def,ghi,jkl")
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	_, err = client.Post.LoadMoreComments(ctx, nil)
	require.EqualError(t, err, "*PostAndComments: cannot be nil")

	resp, err := client.Post.LoadMoreComments(ctx, &PostAndComments{})
	require.Nil(t, resp)
	require.Nil(t, err)

	pc := &PostAndComments{
		Post: &Post{
			FullID: "t3_123",
		},
		Comments: []*Comment{
			{
				FullID: "t1_abc",
			},
		},
		More: &More{
			Children: []string{"def", "ghi", "jkl"},
		},
	}

	_, err = client.Post.LoadMoreComments(ctx, pc)
	require.NoError(t, err)
	require.False(t, pc.HasMore())
	require.Len(t, pc.Comments, 2)
	require.True(t, pc.Comments[1].HasMore())
	require.Len(t, pc.Comments[0].Replies.Comments, 1)
	require.Len(t, pc.Comments[0].Replies.Comments[0].Replies.Comments, 1)
}

func TestPostService_RandomFromSubreddits(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/post/post.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/random", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	postAndComments, _, err := client.Post.RandomFromSubreddits(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, expectedPostAndComments, postAndComments)
}

func TestPostService_Random(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/post/post.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/all/random", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	postAndComments, _, err := client.Post.Random(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedPostAndComments, postAndComments)
}

func TestPostService_RandomFromSubscriptions(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/post/post.json")
	require.NoError(t, err)

	mux.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	postAndComments, _, err := client.Post.RandomFromSubscriptions(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedPostAndComments, postAndComments)
}

func TestPostService_Delete(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/del", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.Delete(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_Save(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/save", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.Save(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_Unsave(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/unsave", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.Unsave(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_EnableReplies(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/sendreplies", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("state", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.EnableReplies(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_DisableReplies(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/sendreplies", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("state", "false")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.DisableReplies(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_Lock(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/lock", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.Lock(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_Unlock(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/unlock", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.Unlock(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_Upvote(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("dir", "1")
		form.Set("rank", "10")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.Upvote(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_Downvote(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("dir", "-1")
		form.Set("rank", "10")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.Downvote(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_RemoveVote(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("dir", "0")
		form.Set("rank", "10")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	resp, err := client.Post.RemoveVote(ctx, "t3_test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostService_MarkVisited(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/store_visits", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("links", "t3_test1,t3_test2,t3_test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Post.MarkVisited(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	_, err = client.Post.MarkVisited(ctx, "t3_test1", "t3_test2", "t3_test3")
	require.NoError(t, err)
}

func TestPostService_Report(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/report", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("thing_id", "t3_test")
		form.Set("reason", "test reason")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Post.Report(ctx, "t3_test", "test reason")
	require.NoError(t, err)
}
