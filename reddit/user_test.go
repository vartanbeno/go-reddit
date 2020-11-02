package reddit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedUser = &User{
	ID:      "test",
	Name:    "Test_User",
	Created: &Timestamp{time.Date(2012, 10, 18, 10, 11, 11, 0, time.UTC)},

	PostKarma:    8239,
	CommentKarma: 130514,

	HasVerifiedEmail: true,
}

var expectedUsers = map[string]*UserSummary{
	"t2_1": {
		Name:         "test_user_1",
		Created:      &Timestamp{time.Date(2017, 3, 12, 2, 1, 47, 0, time.UTC)},
		PostKarma:    488,
		CommentKarma: 22223,
		NSFW:         false,
	},
	"t2_2": {
		Name:         "test_user_2",
		Created:      &Timestamp{time.Date(2015, 12, 20, 18, 12, 51, 0, time.UTC)},
		PostKarma:    8277,
		CommentKarma: 131948,
		NSFW:         false,
	},
	"t2_3": {
		Name:         "test_user_3",
		Created:      &Timestamp{time.Date(2013, 3, 4, 15, 46, 31, 0, time.UTC)},
		PostKarma:    126887,
		CommentKarma: 81918,
		NSFW:         true,
	},
}

var expectedPost = &Post{
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

	SubredditName:         "redditdev",
	SubredditNamePrefixed: "r/redditdev",
	SubredditID:           "t5_2qizd",
	SubredditSubscribers:  37829,

	Author:   "v_95",
	AuthorID: "t2_164ab8",

	IsSelfPost: true,
}

var expectedComment = &Comment{
	ID:      "f0zsa37",
	FullID:  "t1_f0zsa37",
	Created: &Timestamp{time.Date(2019, 9, 21, 21, 38, 16, 0, time.UTC)},
	Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

	ParentID:  "t3_d7ejpn",
	Permalink: "/r/apple/comments/d7ejpn/im_giving_away_an_iphone_11_pro_to_a_commenter_at/f0zsa37/",

	Body:     "Thank you!",
	Author:   "v_95",
	AuthorID: "t2_164ab8",

	SubredditName:         "apple",
	SubredditNamePrefixed: "r/apple",
	SubredditID:           "t5_2qh1f",

	Likes: Bool(true),

	Score:            1,
	Controversiality: 0,

	PostID:          "t3_d7ejpn",
	PostTitle:       "I'm giving away an iPhone 11 Pro to a commenter at random to celebrate Apollo for Reddit's new iOS 13 update and as a thank you to the community! Just leave a comment on this post and the winner will be selected randomly and announced tomorrow at 8 PM GMT. Details inside, and good luck!",
	PostPermalink:   "https://www.reddit.com/r/apple/comments/d7ejpn/im_giving_away_an_iphone_11_pro_to_a_commenter_at/",
	PostAuthor:      "iamthatis",
	PostNumComments: Int(89751),
}

var expectedRelationship = &Relationship{
	ID:      "r9_tqfqp8",
	User:    "test123",
	UserID:  "t2_7b8q1eob",
	Created: &Timestamp{time.Date(2020, 6, 18, 20, 36, 34, 0, time.UTC)},
}

var expectedBlocked = &Blocked{
	Blocked:   "test123",
	BlockedID: "t2_3v9o1yoi",
	Created:   &Timestamp{time.Date(2020, 6, 16, 16, 49, 50, 0, time.UTC)},
}

var expectedTrophies = []*Trophy{
	{
		ID:          "",
		Name:        "Three-Year Club",
		Description: "",
	},
	{
		ID:          "1q1tez",
		Name:        "Verified Email",
		Description: "",
	},
}

var expectedUserSubreddits = []*Subreddit{
	{
		ID:      "3kefx",
		FullID:  "t5_3kefx",
		Created: &Timestamp{time.Date(2017, 5, 11, 16, 37, 16, 0, time.UTC)},

		URL:          "/user/nickofnight/",
		Name:         "u_nickofnight",
		NamePrefixed: "u/nickofnight",
		Title:        "nickofnight",
		Description:  "Stories written for Writing Prompts, NoSleep, and originals. Current series: The Carnival of Night ",
		Type:         "user",
	},
	{
		ID:      "3knn1",
		FullID:  "t5_3knn1",
		Created: &Timestamp{time.Date(2017, 5, 18, 2, 15, 55, 0, time.UTC)},

		URL:                  "/user/shittymorph/",
		Name:                 "u_shittymorph",
		NamePrefixed:         "u/shittymorph",
		Title:                "shittymorph",
		Description:          "In nineteen ninety eight the undertaker threw mankind off hеll in a cell, and plummeted sixteen feet through an announcer's table.",
		Type:                 "user",
		SuggestedCommentSort: "qa",
	},
}

var expectedSearchUsers = []*User{
	{
		ID:      "179965",
		Name:    "washingtonpost",
		Created: &Timestamp{time.Date(2017, 4, 20, 21, 23, 58, 0, time.UTC)},

		PostKarma:    1075227,
		CommentKarma: 339569,

		HasVerifiedEmail: true,
	},
	{
		ID:      "11kowl2w",
		Name:    "reuters",
		Created: &Timestamp{time.Date(2018, 3, 15, 1, 50, 4, 0, time.UTC)},

		PostKarma:    76744,
		CommentKarma: 42717,

		HasVerifiedEmail: true,
	},
}

func TestUserService_Get(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/get.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/Test_User/about", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	user, _, err := client.User.Get(ctx, "Test_User")
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)
}

func TestUserService_GetMultipleByID(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/get-multiple-by-id.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/user_data_by_account_ids", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, "t2_1,t2_2,t2_3", r.Form.Get("ids"))

		fmt.Fprint(w, blob)
	})

	users, _, err := client.User.GetMultipleByID(ctx, "t2_1", "t2_2", "t2_3")
	require.NoError(t, err)
	require.Equal(t, expectedUsers, users)
}

func TestUserService_UsernameAvailable(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/username_available", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)

		user := r.Form.Get("user")
		require.NotEmpty(t, user)

		result := user == "test123"
		fmt.Fprint(w, result)
	})

	ok, _, err := client.User.UsernameAvailable(ctx, "test123")
	require.NoError(t, err)
	require.True(t, ok)

	ok, _, err = client.User.UsernameAvailable(ctx, "123test")
	require.NoError(t, err)
	require.False(t, ok)
}

func TestUserService_Overview(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/overview.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/overview", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, comments, resp, err := client.User.Overview(ctx, nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])

	require.Len(t, comments, 1)
	require.Equal(t, expectedComment, comments[0])

	require.Equal(t, "t1_f0zsa37", resp.After)
}

func TestUserService_OverviewOf(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/overview.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user2/overview", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, comments, resp, err := client.User.OverviewOf(ctx, "user2", nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])

	require.Len(t, comments, 1)
	require.Equal(t, expectedComment, comments[0])

	require.Equal(t, "t1_f0zsa37", resp.After)
}

func TestUserService_Overview_Options(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/overview.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/overview", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "5")
		form.Set("after", "t3_after")
		form.Set("sort", "top")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	_, _, _, err = client.User.Overview(ctx, &ListUserOverviewOptions{
		ListOptions: ListOptions{
			Limit: 5,
			After: "t3_after",
		},
		Sort: "top",
	})
	require.NoError(t, err)
}

func TestUserService_Posts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/submitted.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/submitted", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.User.Posts(ctx, nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])
	require.Equal(t, "t3_gczwql", resp.After)
}

func TestUserService_PostsOf(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/submitted.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user2/submitted", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.User.PostsOf(ctx, "user2", nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])
	require.Equal(t, "t3_gczwql", resp.After)
}

func TestUserService_Posts_Options(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/submitted.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/submitted", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "10")
		form.Set("sort", "new")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	_, _, err = client.User.Posts(ctx, &ListUserOverviewOptions{
		ListOptions: ListOptions{
			Limit: 10,
		},
		Sort: "new",
	})
	require.NoError(t, err)
}

func TestUserService_Comments(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/comments.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/comments", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	comments, resp, err := client.User.Comments(ctx, nil)
	require.NoError(t, err)

	require.Len(t, comments, 1)
	require.Equal(t, expectedComment, comments[0])
	require.Equal(t, "t1_f0zsa37", resp.After)
}

func TestUserService_CommentsOf(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/comments.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user2/comments", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	comments, resp, err := client.User.CommentsOf(ctx, "user2", nil)
	require.NoError(t, err)

	require.Len(t, comments, 1)
	require.Equal(t, expectedComment, comments[0])
	require.Equal(t, "t1_f0zsa37", resp.After)
}

func TestUserService_Comments_Options(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/comments.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/comments", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "100")
		form.Set("before", "t1_before")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	_, _, err = client.User.Comments(ctx, &ListUserOverviewOptions{
		ListOptions: ListOptions{
			Limit:  100,
			Before: "t1_before",
		},
	})
	require.NoError(t, err)
}

func TestUserService_Saved(t *testing.T) {
	client, mux := setup(t)

	// we'll use this, similar payloads
	blob, err := readFileContents("../testdata/user/overview.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/saved", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, comments, resp, err := client.User.Saved(ctx, nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])

	require.Len(t, comments, 1)
	require.Equal(t, expectedComment, comments[0])

	require.Equal(t, "t1_f0zsa37", resp.After)
}

func TestUserService_Saved_Options(t *testing.T) {
	client, mux := setup(t)

	// we'll use this, similar payloads
	blob, err := readFileContents("../testdata/user/overview.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/saved", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "50")
		form.Set("sort", "controversial")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	_, _, _, err = client.User.Saved(ctx, &ListUserOverviewOptions{
		ListOptions: ListOptions{
			Limit: 50,
		},
		Sort: "controversial",
	})
	require.NoError(t, err)
}
func TestUserService_Upvoted(t *testing.T) {
	client, mux := setup(t)

	// we'll use this, similar payloads
	blob, err := readFileContents("../testdata/user/submitted.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/upvoted", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.User.Upvoted(ctx, nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])
	require.Equal(t, "t3_gczwql", resp.After)
}

func TestUserService_Upvoted_Options(t *testing.T) {
	client, mux := setup(t)

	// we'll use this, similar payloads
	blob, err := readFileContents("../testdata/user/submitted.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/upvoted", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "30")
		form.Set("after", "t3_after")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	_, _, err = client.User.Upvoted(ctx, &ListUserOverviewOptions{
		ListOptions: ListOptions{
			Limit: 30,
			After: "t3_after",
		},
	})
	require.NoError(t, err)
}

func TestUserService_UpvotedOf(t *testing.T) {
	client, mux := setup(t)

	// we'll use this, similar payloads
	blob, err := readFileContents("../testdata/user/submitted.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user2/upvoted", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.User.UpvotedOf(ctx, "user2", nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])
	require.Equal(t, "t3_gczwql", resp.After)
}

func TestUserService_Downvoted(t *testing.T) {
	client, mux := setup(t)

	// we'll use this, similar payloads
	blob, err := readFileContents("../testdata/user/submitted.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/downvoted", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.User.Downvoted(ctx, nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])
	require.Equal(t, "t3_gczwql", resp.After)
}

func TestUserService_Downvoted_Options(t *testing.T) {
	client, mux := setup(t)

	// we'll use this, similar payloads
	blob, err := readFileContents("../testdata/user/submitted.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/downvoted", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "20")
		form.Set("before", "t3_before")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	_, _, err = client.User.Downvoted(ctx, &ListUserOverviewOptions{
		ListOptions: ListOptions{
			Limit:  20,
			Before: "t3_before",
		},
	})
	require.NoError(t, err)
}

func TestUserService_DownvotedOf(t *testing.T) {
	client, mux := setup(t)

	// we'll use this, similar payloads
	blob, err := readFileContents("../testdata/user/submitted.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user2/downvoted", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.User.DownvotedOf(ctx, "user2", nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])
	require.Equal(t, "t3_gczwql", resp.After)
}

func TestUserService_Hidden(t *testing.T) {
	client, mux := setup(t)

	// we'll use this, similar payloads
	blob, err := readFileContents("../testdata/user/submitted.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/hidden", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.User.Hidden(ctx, nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])
	require.Equal(t, "t3_gczwql", resp.After)
}

func TestUserService_Gilded(t *testing.T) {
	client, mux := setup(t)

	// we'll use this, similar payloads
	blob, err := readFileContents("../testdata/user/submitted.json")
	require.NoError(t, err)

	mux.HandleFunc("/user/user1/gilded", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.User.Gilded(ctx, nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])
	require.Equal(t, "t3_gczwql", resp.After)
}

func TestUserService_GetFriendship(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/friend.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/v1/me/friends/test123", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	relationship, _, err := client.User.GetFriendship(ctx, "test123")
	require.NoError(t, err)
	require.Equal(t, expectedRelationship, relationship)
}

func TestUserService_Friend(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/friend.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/v1/me/friends/test123", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)

		var request struct {
			Username string `json:"name"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		require.NoError(t, err)
		require.Equal(t, "test123", request.Username)

		fmt.Fprint(w, blob)
	})

	relationship, _, err := client.User.Friend(ctx, "test123")
	require.NoError(t, err)
	require.Equal(t, expectedRelationship, relationship)
}

func TestUserService_Unfriend(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/me/friends/test123", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.User.Unfriend(ctx, "test123")
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUserService_Block(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/block.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/block_user", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("name", "test123")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	blocked, _, err := client.User.Block(ctx, "test123")
	require.NoError(t, err)
	require.Equal(t, expectedBlocked, blocked)
}

func TestUserService_BlockByID(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/block.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/block_user", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("account_id", "abc123")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	blocked, _, err := client.User.BlockByID(ctx, "abc123")
	require.NoError(t, err)
	require.Equal(t, expectedBlocked, blocked)
}

func TestUserService_Unblock(t *testing.T) {
	client, mux := setup(t)

	client.redditID = "self123"

	mux.HandleFunc("/api/unfriend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("name", "test123")
		form.Set("type", "enemy")
		form.Set("container", client.redditID)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.User.Unblock(ctx, "test123")
	require.NoError(t, err)
}

func TestUserService_UnblockByID(t *testing.T) {
	client, mux := setup(t)

	client.redditID = "self123"

	mux.HandleFunc("/api/unfriend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "abc123")
		form.Set("type", "enemy")
		form.Set("container", client.redditID)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.User.UnblockByID(ctx, "abc123")
	require.NoError(t, err)
}

func TestUserService_Trophies(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/trophies.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/v1/user/user1/trophies", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	trophies, _, err := client.User.Trophies(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedTrophies, trophies)
}

func TestUserService_TrophiesOf(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/trophies.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/v1/user/test123/trophies", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	trophies, _, err := client.User.TrophiesOf(ctx, "test123")
	require.NoError(t, err)
	require.Equal(t, expectedTrophies, trophies)
}

func TestUserService_Popular(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/user-subreddits.json")
	require.NoError(t, err)

	mux.HandleFunc("/users/popular", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	userSubreddits, resp, err := client.User.Popular(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedUserSubreddits, userSubreddits)
	require.Equal(t, "t5_3knn1", resp.After)
}

func TestUserService_New(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/user-subreddits.json")
	require.NoError(t, err)

	mux.HandleFunc("/users/new", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	userSubreddits, resp, err := client.User.New(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedUserSubreddits, userSubreddits)
	require.Equal(t, "t5_3knn1", resp.After)
}

func TestUserService_Search(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/user/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/users/search", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("q", "test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	users, resp, err := client.User.Search(ctx, "test", nil)
	require.NoError(t, err)
	require.Equal(t, expectedSearchUsers, users)
	require.Equal(t, "t2_11kowl2w", resp.After)
}
