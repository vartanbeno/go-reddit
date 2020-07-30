package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var expectedPostAndComments = &PostAndComments{
	Post: &Post{
		ID:      "testpost",
		FullID:  "t3_testpost",
		Created: &Timestamp{time.Date(2020, 7, 18, 10, 26, 7, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: Permalink("https://www.reddit.com/r/test/comments/testpost/test/"),
		URL:       "https://www.reddit.com/r/test/comments/testpost/test/",

		Title: "Test",
		Body:  "Hello",

		Score:            1,
		UpvoteRatio:      1,
		NumberOfComments: 2,

		SubredditID:           "t5_2qh23",
		SubredditName:         "test",
		SubredditNamePrefixed: "r/test",

		AuthorID:   "t2_testuser",
		AuthorName: "testuser",

		IsSelfPost: true,
	},
	Comments: []*Comment{
		{
			ID:      "testc1",
			FullID:  "t1_testc1",
			Created: &Timestamp{time.Date(2020, 7, 18, 10, 31, 59, 0, time.UTC)},
			Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

			ParentID:  "t3_testpost",
			Permalink: Permalink("https://www.reddit.com/r/test/comments/testpost/test/testc1/"),

			Body:     "Hi",
			Author:   "testuser",
			AuthorID: "t2_testuser",

			Subreddit:             "test",
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
						Permalink: Permalink("https://www.reddit.com/r/test/comments/testpost/test/testc2/"),

						Body:     "Hello",
						Author:   "testuser",
						AuthorID: "t2_testuser",

						Subreddit:             "test",
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

	Permalink: "https://www.reddit.com/r/test/comments/hw6l6a/test_title/",
	URL:       "https://www.reddit.com/r/test/comments/hw6l6a/test_title/",

	Title: "Test Title",
	Body:  "this is edited",

	Likes: Bool(true),

	Score:            1,
	UpvoteRatio:      1,
	NumberOfComments: 0,

	SubredditID:           "t5_2qh23",
	SubredditName:         "test",
	SubredditNamePrefixed: "r/test",

	AuthorID:   "t2_164ab8",
	AuthorName: "v_95",

	Spoiler:    true,
	IsSelfPost: true,
}

func TestPostService_Get(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/post/post.json")
	assert.NoError(t, err)

	mux.HandleFunc("/comments/test", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	postAndComments, _, err := client.Post.Get(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, expectedPostAndComments, postAndComments)
}

func TestPostService_SubmitText(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/post/submit.json")
	assert.NoError(t, err)

	mux.HandleFunc("/api/submit", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("kind", "self")
		form.Set("sr", "test")
		form.Set("title", "Test Title")
		form.Set("text", "Test Text")
		form.Set("spoiler", "true")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	submittedPost, _, err := client.Post.SubmitText(ctx, SubmitTextOptions{
		Subreddit: "test",
		Title:     "Test Title",
		Text:      "Test Text",
		Spoiler:   true,
	})
	assert.NoError(t, err)
	assert.Equal(t, expectedSubmittedPost, submittedPost)
}

func TestPostService_SubmitLink(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/post/submit.json")
	assert.NoError(t, err)

	mux.HandleFunc("/api/submit", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

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
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	submittedPost, _, err := client.Post.SubmitLink(ctx, SubmitLinkOptions{
		Subreddit:   "test",
		Title:       "Test Title",
		URL:         "https://www.example.com",
		SendReplies: Bool(false),
		Resubmit:    true,
		NSFW:        true,
	})
	assert.NoError(t, err)
	assert.Equal(t, expectedSubmittedPost, submittedPost)
}

func TestPostService_Edit(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/post/edit.json")
	assert.NoError(t, err)

	mux.HandleFunc("/api/editusertext", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("return_rtjson", "true")
		form.Set("thing_id", "t3_test")
		form.Set("text", "test edit")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	editedPost, _, err := client.Post.Edit(ctx, "t3_test", "test edit")
	assert.NoError(t, err)
	assert.Equal(t, expectedEditedPost, editedPost)
}

func TestPostService_Hide(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/hide", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "1,2,3")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	_, err := client.Post.Hide(ctx)
	assert.EqualError(t, err, "must provide at least 1 id")

	res, err := client.Post.Hide(ctx, "1", "2", "3")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Unhide(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unhide", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "1,2,3")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	_, err := client.Post.Unhide(ctx)
	assert.EqualError(t, err, "must provide at least 1 id")

	res, err := client.Post.Unhide(ctx, "1", "2", "3")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_MarkNSFW(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/marknsfw", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.MarkNSFW(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_UnmarkNSFW(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unmarknsfw", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.UnmarkNSFW(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Spoiler(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/spoiler", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Spoiler(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Unspoiler(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unspoiler", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Unspoiler(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Sticky(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_subreddit_sticky", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "true")
		form.Set("num", "1")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Sticky(ctx, "t3_test", false)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Unsticky(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_subreddit_sticky", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "false")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Unsticky(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_PinToProfile(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_subreddit_sticky", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "true")
		form.Set("to_profile", "true")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.PinToProfile(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_UnpinFromProfile(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_subreddit_sticky", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "false")
		form.Set("to_profile", "true")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.UnpinFromProfile(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_SetSuggestedSortBest(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "confidence")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.SetSuggestedSortBest(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_SetSuggestedSortTop(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "top")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.SetSuggestedSortTop(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_SetSuggestedSortNew(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "new")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.SetSuggestedSortNew(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_SetSuggestedSortControversial(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "controversial")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.SetSuggestedSortControversial(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_SetSuggestedSortOld(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "old")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.SetSuggestedSortOld(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_SetSuggestedSortRandom(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "random")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.SetSuggestedSortRandom(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_SetSuggestedSortAMA(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "qa")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.SetSuggestedSortAMA(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_SetSuggestedSortLive(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "live")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.SetSuggestedSortLive(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_ClearSuggestedSort(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_suggested_sort", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("sort", "")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.ClearSuggestedSort(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_EnableContestMode(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_contest_mode", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "true")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.EnableContestMode(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_DisableContestMode(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/set_contest_mode", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t3_test")
		form.Set("state", "false")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.DisableContestMode(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_More(t *testing.T) {
	setup()
	defer teardown()

	parentComment := &Comment{
		FullID:   "t1_abc",
		ParentID: "t3_123",
		PostID:   "t3_123",
		Replies: Replies{
			MoreComments: &More{
				Children: []string{"def,ghi"},
			},
		},
	}

	blob, err := readFileContents("testdata/post/more.json")
	assert.NoError(t, err)

	mux.HandleFunc("/api/morechildren", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("link_id", "t3_123")
		form.Set("children", "def,ghi")
		form.Set("api_type", "json")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	_, err = client.Comment.LoadMoreReplies(ctx, parentComment)
	assert.NoError(t, err)
	assert.Nil(t, parentComment.Replies.MoreComments)
	assert.Len(t, parentComment.Replies.Comments, 1)
	assert.Len(t, parentComment.Replies.Comments[0].Replies.Comments, 1)
}

func TestPostService_MoreNil(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Comment.LoadMoreReplies(ctx, nil)
	assert.EqualError(t, err, "comment: must not be nil")

	parentComment := &Comment{
		Replies: Replies{
			MoreComments: nil,
		},
	}

	// should return nil, nil since comment does not have More struct
	resp, err := client.Comment.LoadMoreReplies(ctx, parentComment)
	assert.NoError(t, err)
	assert.Nil(t, resp)

	parentComment.Replies.MoreComments = &More{
		Children: []string{},
	}

	// should return nil, nil since comment's More struct has 0 children
	resp, err = client.Comment.LoadMoreReplies(ctx, parentComment)
	assert.NoError(t, err)
	assert.Nil(t, resp)
}

func TestPostService_RandomFromSubreddits(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/post/post.json")
	assert.NoError(t, err)

	mux.HandleFunc("/r/test/random", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	postAndComments, _, err := client.Post.RandomFromSubreddits(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, expectedPostAndComments, postAndComments)
}

func TestPostService_Random(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/post/post.json")
	assert.NoError(t, err)

	mux.HandleFunc("/r/all/random", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	postAndComments, _, err := client.Post.Random(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedPostAndComments, postAndComments)
}

func TestPostService_RandomFromSubscriptions(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("testdata/post/post.json")
	assert.NoError(t, err)

	mux.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	postAndComments, _, err := client.Post.RandomFromSubscriptions(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedPostAndComments, postAndComments)
}

func TestPostService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/del", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Delete(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Save(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/save", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Save(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Unsave(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unsave", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Unsave(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_EnableReplies(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/sendreplies", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("state", "true")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.EnableReplies(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_DisableReplies(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/sendreplies", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("state", "false")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.DisableReplies(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Lock(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/lock", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Lock(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Unlock(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unlock", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Unlock(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Upvote(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("dir", "1")
		form.Set("rank", "10")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Upvote(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_Downvote(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("dir", "-1")
		form.Set("rank", "10")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.Downvote(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService_RemoveVote(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("dir", "0")
		form.Set("rank", "10")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.PostForm)
	})

	res, err := client.Post.RemoveVote(ctx, "t3_test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
