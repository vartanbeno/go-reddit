package geddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var expectedUser = &User{
	ID:      "test",
	Name:    "Test_User",
	Created: &Timestamp{time.Date(2012, 10, 18, 10, 11, 11, 0, time.UTC)},

	LinkKarma:    8239,
	CommentKarma: 130514,

	HasVerifiedEmail: true,
}

var expectedUsers = map[string]*UserShort{
	"t2_1": {
		Name:         "test_user_1",
		Created:      &Timestamp{time.Date(2017, 3, 12, 2, 1, 47, 0, time.UTC)},
		LinkKarma:    488,
		CommentKarma: 22223,
		NSFW:         false,
	},
	"t2_2": {
		Name:         "test_user_2",
		Created:      &Timestamp{time.Date(2015, 12, 20, 18, 12, 51, 0, time.UTC)},
		LinkKarma:    8277,
		CommentKarma: 131948,
		NSFW:         false,
	},
	"t2_3": {
		Name:         "test_user_3",
		Created:      &Timestamp{time.Date(2013, 3, 4, 15, 46, 31, 0, time.UTC)},
		LinkKarma:    126887,
		CommentKarma: 81918,
		NSFW:         true,
	},
}

var expectedPost = Link{
	ID:      "gczwql",
	FullID:  "t3_gczwql",
	Created: &Timestamp{time.Date(2020, 5, 3, 22, 46, 25, 0, time.UTC)},
	Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

	Permalink: "/r/redditdev/comments/gczwql/get_userusernamegilded_does_it_return_other_users/",
	URL:       "https://www.reddit.com/r/redditdev/comments/gczwql/get_userusernamegilded_does_it_return_other_users/",

	Title: "GET /user/{username}/gilded: does it return other user's things you've gilded, or your things that have been gilded? Does it return both comments and posts?",
	Body:  "Talking about [this](https://www.reddit.com/dev/api/#GET_user_{username}_{where}) endpoint specifically.\n\nI'm building a Reddit API client, but don't have gold.",

	Likes: Bool(true),

	Score:            9,
	UpvoteRatio:      0.86,
	NumberOfComments: 2,

	SubredditID:           "t5_2qizd",
	SubredditName:         "redditdev",
	SubredditNamePrefixed: "r/redditdev",

	AuthorID:   "t2_164ab8",
	AuthorName: "v_95",

	IsSelfPost: true,
}

var expectedComment = Comment{
	ID:      "f0zsa37",
	FullID:  "t1_f0zsa37",
	Created: &Timestamp{time.Date(2019, 9, 21, 21, 38, 16, 0, time.UTC)},
	Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

	ParentID:  "t3_d7ejpn",
	Permalink: "/r/apple/comments/d7ejpn/im_giving_away_an_iphone_11_pro_to_a_commenter_at/f0zsa37/",

	Body:     "Thank you!",
	Author:   "v_95",
	AuthorID: "t2_164ab8",

	Subreddit:             "apple",
	SubredditNamePrefixed: "r/apple",
	SubredditID:           "t5_2qh1f",

	Likes: Bool(true),

	Score:            1,
	Controversiality: 0,

	LinkID:          "t3_d7ejpn",
	LinkTitle:       "I'm giving away an iPhone 11 Pro to a commenter at random to celebrate Apollo for Reddit's new iOS 13 update and as a thank you to the community! Just leave a comment on this post and the winner will be selected randomly and announced tomorrow at 8 PM GMT. Details inside, and good luck!",
	LinkPermalink:   "https://www.reddit.com/r/apple/comments/d7ejpn/im_giving_away_an_iphone_11_pro_to_a_commenter_at/",
	LinkAuthor:      "iamthatis",
	LinkNumComments: 89751,
}

func TestUserService_Get(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/user/get.json")

	mux.HandleFunc("/user/Test_User/about", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	user, _, err := client.User.Get(ctx, "Test_User")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserService_GetMultipleByID(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/user/get-multiple-by-id.json")

	mux.HandleFunc("/api/user_data_by_account_ids", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "t2_1,t2_2,t2_3", r.Form.Get("ids"))

		fmt.Fprint(w, blob)
	})

	users, _, err := client.User.GetMultipleByID(ctx, "t2_1", "t2_2", "t2_3")
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}

func TestUserService_UsernameAvailable(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/username_available", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		assert.NoError(t, err)

		user := r.Form.Get("user")
		assert.NotEmpty(t, user)

		result := user == "test123"
		fmt.Fprint(w, result)
	})

	ok, _, err := client.User.UsernameAvailable(ctx, "test123")
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, _, err = client.User.UsernameAvailable(ctx, "123test")
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestUserService_Overview(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/user/overview.json")

	mux.HandleFunc("/user/user1/overview", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, comments, _, err := client.User.Overview(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, posts)
	assert.NotNil(t, comments)

	assert.Len(t, posts.Links, 1)
	assert.Equal(t, expectedPost, posts.Links[0])
	assert.Equal(t, "t1_f0zsa37", posts.After)
	assert.Equal(t, "", posts.Before)

	assert.Len(t, comments.Comments, 1)
	assert.Equal(t, expectedComment, comments.Comments[0])
	assert.Equal(t, "t1_f0zsa37", comments.After)
	assert.Equal(t, "", comments.Before)
}

func TestUserService_OverviewOf(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/user/overview.json")

	mux.HandleFunc("/user/user2/overview", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, comments, _, err := client.User.OverviewOf(ctx, "user2")
	assert.NoError(t, err)
	assert.NotNil(t, posts)
	assert.NotNil(t, comments)

	assert.Len(t, posts.Links, 1)
	assert.Equal(t, expectedPost, posts.Links[0])
	assert.Equal(t, "t1_f0zsa37", posts.After)
	assert.Equal(t, "", posts.Before)

	assert.Len(t, comments.Comments, 1)
	assert.Equal(t, expectedComment, comments.Comments[0])
	assert.Equal(t, "t1_f0zsa37", comments.After)
	assert.Equal(t, "", comments.Before)
}

func TestUserService_Overview_Options(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/user/overview.json")

	mux.HandleFunc("/user/user1/overview", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "5")
		form.Set("after", "t3_after")
		form.Set("sort", SortTop.String())

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	_, _, _, err := client.User.Overview(ctx, SetLimit(5), SetAfter("t3_after"), SetSort(SortTop))
	assert.NoError(t, err)
}

func TestUserService_Posts(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/user/submitted.json")

	mux.HandleFunc("/user/user1/submitted", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, _, err := client.User.Posts(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, posts)

	assert.Len(t, posts.Links, 1)
	assert.Equal(t, expectedPost, posts.Links[0])
	assert.Equal(t, "t3_gczwql", posts.After)
	assert.Equal(t, "", posts.Before)
}

func TestUserService_PostsOf(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/user/submitted.json")

	mux.HandleFunc("/user/user2/submitted", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, _, err := client.User.PostsOf(ctx, "user2")
	assert.NoError(t, err)
	assert.NotNil(t, posts)

	assert.Len(t, posts.Links, 1)
	assert.Equal(t, expectedPost, posts.Links[0])
	assert.Equal(t, "t3_gczwql", posts.After)
	assert.Equal(t, "", posts.Before)
}

func TestUserService_Posts_Options(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/user/submitted.json")

	mux.HandleFunc("/user/user1/submitted", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "10")
		form.Set("sort", SortNew.String())

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	_, _, err := client.User.Posts(ctx, SetLimit(10), SetSort(SortNew))
	assert.NoError(t, err)
}

func TestUserService_Comments(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/user/comments.json")

	mux.HandleFunc("/user/user1/comments", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	comments, _, err := client.User.Comments(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, comments)

	assert.Len(t, comments.Comments, 1)
	assert.Equal(t, expectedComment, comments.Comments[0])
	assert.Equal(t, "t1_f0zsa37", comments.After)
	assert.Equal(t, "", comments.Before)
}

func TestUserService_CommentsOf(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/user/comments.json")

	mux.HandleFunc("/user/user2/comments", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	comments, _, err := client.User.CommentsOf(ctx, "user2")
	assert.NoError(t, err)
	assert.NotNil(t, comments)

	assert.Len(t, comments.Comments, 1)
	assert.Equal(t, expectedComment, comments.Comments[0])
	assert.Equal(t, "t1_f0zsa37", comments.After)
	assert.Equal(t, "", comments.Before)
}

func TestUserService_Comments_Options(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/user/comments.json")

	mux.HandleFunc("/user/user1/comments", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "100")
		form.Set("before", "t1_before")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	_, _, err := client.User.Comments(ctx, SetLimit(100), SetBefore("t1_before"))
	assert.NoError(t, err)
}

func TestUserService_Saved(t *testing.T) {
	setup()
	defer teardown()

	// we'll use this, similar payloads
	blob := readFileContents(t, "testdata/user/overview.json")

	mux.HandleFunc("/user/user1/saved", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, comments, _, err := client.User.Saved(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, posts)
	assert.NotNil(t, comments)

	assert.Len(t, posts.Links, 1)
	assert.Equal(t, expectedPost, posts.Links[0])
	assert.Equal(t, "t1_f0zsa37", posts.After)
	assert.Equal(t, "", posts.Before)

	assert.Len(t, comments.Comments, 1)
	assert.Equal(t, expectedComment, comments.Comments[0])
	assert.Equal(t, "t1_f0zsa37", comments.After)
	assert.Equal(t, "", comments.Before)
}

func TestUserService_Saved_Options(t *testing.T) {
	setup()
	defer teardown()

	// we'll use this, similar payloads
	blob := readFileContents(t, "testdata/user/overview.json")

	mux.HandleFunc("/user/user1/saved", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "50")
		form.Set("sort", SortControversial.String())

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	_, _, _, err := client.User.Saved(ctx, SetLimit(50), SetSort(SortControversial))
	assert.NoError(t, err)
}
func TestUserService_Upvoted(t *testing.T) {
	setup()
	defer teardown()

	// we'll use this, similar payloads
	blob := readFileContents(t, "testdata/user/submitted.json")

	mux.HandleFunc("/user/user1/upvoted", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, _, err := client.User.Upvoted(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, posts)

	assert.Len(t, posts.Links, 1)
	assert.Equal(t, expectedPost, posts.Links[0])
	assert.Equal(t, "t3_gczwql", posts.After)
	assert.Equal(t, "", posts.Before)
}

func TestUserService_Upvoted_Options(t *testing.T) {
	setup()
	defer teardown()

	// we'll use this, similar payloads
	blob := readFileContents(t, "testdata/user/submitted.json")

	mux.HandleFunc("/user/user1/upvoted", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "30")
		form.Set("after", "t3_after")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	_, _, err := client.User.Upvoted(ctx, SetLimit(30), SetAfter("t3_after"))
	assert.NoError(t, err)
}
func TestUserService_Downvoted(t *testing.T) {
	setup()
	defer teardown()

	// we'll use this, similar payloads
	blob := readFileContents(t, "testdata/user/submitted.json")

	mux.HandleFunc("/user/user1/downvoted", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, _, err := client.User.Downvoted(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, posts)

	assert.Len(t, posts.Links, 1)
	assert.Equal(t, expectedPost, posts.Links[0])
	assert.Equal(t, "t3_gczwql", posts.After)
	assert.Equal(t, "", posts.Before)
}

func TestUserService_Downvoted_Options(t *testing.T) {
	setup()
	defer teardown()

	// we'll use this, similar payloads
	blob := readFileContents(t, "testdata/user/submitted.json")

	mux.HandleFunc("/user/user1/downvoted", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "20")
		form.Set("before", "t3_before")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	_, _, err := client.User.Downvoted(ctx, SetLimit(20), SetBefore("t3_before"))
	assert.NoError(t, err)
}

func TestUserService_Hidden(t *testing.T) {
	setup()
	defer teardown()

	// we'll use this, similar payloads
	blob := readFileContents(t, "testdata/user/submitted.json")

	mux.HandleFunc("/user/user1/hidden", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, _, err := client.User.Hidden(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, posts)

	assert.Len(t, posts.Links, 1)
	assert.Equal(t, expectedPost, posts.Links[0])
	assert.Equal(t, "t3_gczwql", posts.After)
	assert.Equal(t, "", posts.Before)
}

func TestUserService_Gilded(t *testing.T) {
	setup()
	defer teardown()

	// we'll use this, similar payloads
	blob := readFileContents(t, "testdata/user/submitted.json")

	mux.HandleFunc("/user/user1/gilded", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, _, err := client.User.Gilded(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, posts)

	assert.Len(t, posts.Links, 1)
	assert.Equal(t, expectedPost, posts.Links[0])
	assert.Equal(t, "t3_gczwql", posts.After)
	assert.Equal(t, "", posts.Before)
}
