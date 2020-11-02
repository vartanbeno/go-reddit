package reddit

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedPosts = []*Post{
	{
		ID:      "agi5zf",
		FullID:  "t3_agi5zf",
		Created: &Timestamp{time.Date(2019, 1, 16, 5, 57, 51, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/test/comments/agi5zf/test/",
		URL:       "https://www.reddit.com/r/test/comments/agi5zf/test/",

		Title: "test",
		Body:  "test",

		Score:            253,
		UpvoteRatio:      0.99,
		NumberOfComments: 1634,

		SubredditName:         "test",
		SubredditNamePrefixed: "r/test",
		SubredditID:           "t5_2qh23",
		SubredditSubscribers:  8154,

		Author:   "kmiller0112",
		AuthorID: "t2_30a5ktgt",

		IsSelfPost: true,
		Stickied:   true,
	},
	{
		ID:      "hyhquk",
		FullID:  "t3_hyhquk",
		Created: &Timestamp{time.Date(2020, 7, 27, 0, 5, 10, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/test/comments/hyhquk/veggies/",
		URL:       "https://i.imgur.com/LrN2mPw.jpg",

		Title: "Veggies",

		Score:            4,
		UpvoteRatio:      1,
		NumberOfComments: 0,

		SubredditName:         "test",
		SubredditNamePrefixed: "r/test",
		SubredditID:           "t5_2qh23",
		SubredditSubscribers:  8154,

		Author:   "MuckleMcDuckle",
		AuthorID: "t2_6fqntbwq",
	},
}

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
	Subscribed:      true,
}

var expectedSubreddits = []*Subreddit{
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
		Subscribed:  true,
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
		Subscribed:  true,
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
		Subscribed:  false,
		Favorite:    false,
	},
}

var expectedSubredditNames = []string{
	"golang",
	"golang_infosec",
	"GolangJobOfferings",
	"golanguage",
	"golang2",
	"GolangUnofficial",
	"golang_jobs",
}

var expectedSearchPosts = []*Post{
	{
		ID:      "hybow9",
		FullID:  "t3_hybow9",
		Created: &Timestamp{time.Date(2020, 7, 26, 18, 14, 24, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/WatchPeopleDieInside/comments/hybow9/pregnancy_test/",
		URL:       "https://v.redd.it/ra4qnt8bt8d51",

		Title: "Pregnancy test",

		Score:            103829,
		UpvoteRatio:      0.88,
		NumberOfComments: 3748,

		SubredditName:         "WatchPeopleDieInside",
		SubredditNamePrefixed: "r/WatchPeopleDieInside",
		SubredditID:           "t5_3h4zq",
		SubredditSubscribers:  2599948,

		Author:   "chocolat_ice_cream",
		AuthorID: "t2_3p32m02",
	},
	{
		ID:      "hmwhd7",
		FullID:  "t3_hmwhd7",
		Created: &Timestamp{time.Date(2020, 7, 7, 15, 19, 42, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/worldnews/comments/hmwhd7/brazilian_president_jair_bolsonaro_tests_positive/",
		URL:       "https://www.theguardian.com/world/2020/jul/07/jair-bolsonaro-coronavirus-positive-test-brazil-president",

		Title: "Brazilian president Jair Bolsonaro tests positive for coronavirus",

		Score:            149238,
		UpvoteRatio:      0.94,
		NumberOfComments: 7415,

		SubredditName:         "worldnews",
		SubredditNamePrefixed: "r/worldnews",
		SubredditID:           "t5_2qh13",
		SubredditSubscribers:  24651441,

		Author:   "Jeremy_Martin",
		AuthorID: "t2_wgrkg",
	},
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

var expectedRelationships3 = []*Relationship{
	{
		ID:      "rel_id1",
		Created: &Timestamp{time.Date(2020, 8, 11, 2, 35, 2, 0, time.UTC)},
		User:    "testuser1",
		UserID:  "t2_user1",
	},
	{
		ID:      "rel_id2",
		Created: &Timestamp{time.Date(2020, 8, 11, 2, 35, 0, 0, time.UTC)},
		User:    "testuser2",
		UserID:  "t2_user2",
	},
}

var expectedBans = []*Ban{
	{
		Relationship: &Relationship{
			ID:      "rb_123",
			Created: &Timestamp{time.Date(2020, 8, 11, 2, 35, 2, 0, time.UTC)},

			User:   "testuser1",
			UserID: "t2_user1",
		},

		DaysLeft: Int(43),
		Note:     "Spam",
	},
	{
		Relationship: &Relationship{
			ID:      "rb_456",
			Created: &Timestamp{time.Date(2020, 8, 11, 2, 35, 0, 0, time.UTC)},

			User:   "testuser2",
			UserID: "t2_user2",
		},

		DaysLeft: nil,
		Note:     "Spam",
	},
}

var expectedModerators = []*Moderator{
	{
		Relationship: &Relationship{
			ID:      "rb_tmatb9",
			User:    "testuser1",
			UserID:  "t2_test1",
			Created: &Timestamp{time.Date(2013, 7, 29, 20, 44, 27, 0, time.UTC)},
		},
		Permissions: []string{"all"},
	},
	{
		Relationship: &Relationship{
			ID:      "rb_5c9s4d",
			User:    "testuser2",
			UserID:  "t2_test2",
			Created: &Timestamp{time.Date(2014, 3, 1, 18, 13, 53, 0, time.UTC)},
		},
		Permissions: []string{"all"},
	},
}

var expectedRules = []*SubredditRule{
	{
		Kind:            "link",
		Name:            "Read the Rules Before Posting",
		ViolationReason: "Read the Rules Before Posting",
		Description:     "https://www.reddit.com/r/Fitness/wiki/rules",
		Priority:        0,
		Created:         &Timestamp{time.Date(2019, 5, 22, 5, 32, 58, 0, time.UTC)},
	},
	{
		Kind:            "link",
		Name:            "Read the Wiki Before Posting",
		ViolationReason: "Read the Wiki Before Posting",
		Description:     "https://thefitness.wiki",
		Priority:        1,
		Created:         &Timestamp{time.Date(2019, 11, 9, 7, 56, 33, 0, time.UTC)},
	},
}

var expectedDayTraffic = []*SubredditTrafficStats{
	{&Timestamp{time.Date(2020, 9, 13, 0, 0, 0, 0, time.UTC)}, 0, 0, 0},
	{&Timestamp{time.Date(2020, 9, 12, 0, 0, 0, 0, time.UTC)}, 1, 12, 0},
	{&Timestamp{time.Date(2020, 9, 11, 0, 0, 0, 0, time.UTC)}, 5, 85, 0},
	{&Timestamp{time.Date(2020, 9, 10, 0, 0, 0, 0, time.UTC)}, 4, 20, 0},
	{&Timestamp{time.Date(2020, 9, 9, 0, 0, 0, 0, time.UTC)}, 2, 64, 0},
	{&Timestamp{time.Date(2020, 9, 8, 0, 0, 0, 0, time.UTC)}, 2, 95, 0},
	{&Timestamp{time.Date(2020, 9, 7, 0, 0, 0, 0, time.UTC)}, 3, 41, 0},
}

var expectedHourTraffic = []*SubredditTrafficStats{
	{&Timestamp{time.Date(2020, 9, 12, 20, 0, 0, 0, time.UTC)}, 1, 12, 0},
	{&Timestamp{time.Date(2020, 9, 11, 3, 0, 0, 0, time.UTC)}, 4, 57, 0},
	{&Timestamp{time.Date(2020, 9, 11, 2, 0, 0, 0, time.UTC)}, 4, 28, 0},
}

var expectedMonthTraffic = []*SubredditTrafficStats{
	{&Timestamp{time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC)}, 7, 481, 0},
	{&Timestamp{time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC)}, 5, 346, 0},
	{&Timestamp{time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)}, 4, 264, 0},
}

var expectedStyleSheet = &SubredditStyleSheet{
	SubredditID: "t5_2rc7j",
	Images: []*SubredditImage{
		{
			Name: "gopher",
			Link: "url(%%gopher%%)",
			URL:  "http://b.thumbs.redditmedia.com/q5Wb6hTPm2Bd6Of9_xMrTu4n5qgAljJNqtnbE3Tging.png",
		},
	},
	StyleSheet: `.flair-gopher {
    background: url(%%gopher%%) no-repeat;
    border: 0;
    padding: 0;
    width: 16px;
    height: 16px;
}`,
}

var expectedSubredditSettings = &SubredditSettings{
	ID: "t5_test",

	Type: String("private"),

	Language: String("en"),

	Title:                 String("hello!"),
	Description:           String("description"),
	Sidebar:               String("sidebar"),
	SubmissionText:        String(""),
	WelcomeMessage:        String(""),
	WelcomeMessageEnabled: Bool(false),

	AllowCrossposts:            Bool(false),
	AllowChatPosts:             Bool(true),
	AllowPollPosts:             Bool(false),
	AllowFreeFormReports:       Bool(true),
	AllowOriginalContent:       Bool(false),
	AllowImages:                Bool(true),
	AllowMultipleImagesPerPost: Bool(true),

	ExcludeSitewideBannedUsersContent: Bool(false),

	CrowdControlChalLevel: Int(2),

	AllOriginalContent: Bool(false),

	SuggestedCommentSort: nil,

	SubmitLinkPostLabel: String("submit a link!"),
	SubmitTextPostLabel: String("submit a post!"),

	PostType: String("any"),

	SpamFilterStrengthLinkPosts: String("low"),
	SpamFilterStrengthTextPosts: String("low"),
	SpamFilterStrengthComments:  String("low"),

	ShowContentThumbnails:              Bool(false),
	ExpandMediaPreviewsOnCommentsPages: Bool(true),

	CollapseDeletedComments:    Bool(false),
	MinutesToHideCommentScores: Int(0),

	SpoilersEnabled: Bool(true),

	HeaderMouseoverText: String("hello!"),

	MobileColour: String(""),

	HideAds: Bool(false),
	NSFW:    Bool(false),

	AllowDiscoveryInHighTrafficFeeds: Bool(true),
	AllowDiscoveryByIndividualUsers:  Bool(true),

	WikiMode:              String("modonly"),
	WikiMinimumAccountAge: Int(0),
	WikiMinimumKarma:      Int(0),
}

var expectedSubredditPostRequirements = &SubredditPostRequirements{
	Guidelines:              "test",
	GuidelinesDisplayPolicy: "",

	TitleMinLength: 50,
	TitleMaxLength: 200,

	BodyMinLength: 50,
	BodyMaxLength: 2000,

	TitleBlacklistedStrings: []string{"no"},
	BodyBlacklistedStrings:  []string{"no"},

	TitleRequiredStrings: []string{"yes"},
	BodyRequiredStrings:  []string{"yes"},

	DomainBlacklist: []string{"example.com"},
	DomainWhitelist: []string{},

	BodyRestrictionPolicy: "none",
	LinkRestrictionPolicy: "none",

	GalleryMinItems:            2,
	GalleryMaxItems:            20,
	GalleryCaptionsRequirement: "none",
	GalleryURLsRequirement:     "none",

	LinkRepostAge: 2,
	FlairRequired: false,

	TitleRegexes: []string{},
	BodyRegexes:  []string{},
}

func TestSubredditService_HotPosts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/hot", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.HotPosts(ctx, "test", nil)
	require.NoError(t, err)
	require.Equal(t, expectedPosts, posts)
	require.Equal(t, "t3_hyhquk", resp.After)
}

func TestSubredditService_NewPosts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/new", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.NewPosts(ctx, "test", nil)
	require.NoError(t, err)
	require.Equal(t, expectedPosts, posts)
	require.Equal(t, "t3_hyhquk", resp.After)
}

func TestSubredditService_RisingPosts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/rising", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.RisingPosts(ctx, "test", nil)
	require.NoError(t, err)
	require.Equal(t, expectedPosts, posts)
	require.Equal(t, "t3_hyhquk", resp.After)
}

func TestSubredditService_ControversialPosts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/controversial", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.ControversialPosts(ctx, "test", nil)
	require.NoError(t, err)
	require.Equal(t, expectedPosts, posts)
	require.Equal(t, "t3_hyhquk", resp.After)
}

func TestSubredditService_TopPosts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/top", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.TopPosts(ctx, "test", nil)
	require.NoError(t, err)
	require.Equal(t, expectedPosts, posts)
	require.Equal(t, "t3_hyhquk", resp.After)
}

func TestSubredditService_Get(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/about.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/golang/about", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	_, _, err = client.Subreddit.Get(ctx, "")
	require.EqualError(t, err, "name: cannot be empty")

	subreddit, _, err := client.Subreddit.Get(ctx, "golang")
	require.NoError(t, err)
	require.Equal(t, expectedSubreddit, subreddit)
}

func TestSubredditService_Popular(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/popular", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Popular(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_New(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/new", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.New(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_Gold(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/gold", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Gold(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_Default(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/default", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Default(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_Subscribed(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/mine/subscriber", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Subscribed(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_Approved(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/mine/contributor", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Approved(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_Moderated(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/mine/moderator", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Moderated(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_GetSticky1(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/post/post.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/sticky", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, "1", r.Form.Get("num"))

		fmt.Fprint(w, blob)
	})

	postAndComments, _, err := client.Subreddit.GetSticky1(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, expectedPostAndComments, postAndComments)
}

func TestSubredditService_GetSticky2(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/post/post.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/sticky", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, "2", r.Form.Get("num"))

		fmt.Fprint(w, blob)
	})

	postAndComments, _, err := client.Subreddit.GetSticky2(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, expectedPostAndComments, postAndComments)
}

func TestSubredditService_Subscribe(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("action", "sub")
		form.Set("sr_name", "test,golang,nba")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.Subscribe(ctx, "test", "golang", "nba")
	require.NoError(t, err)
}

func TestSubredditService_SubscribeByID(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("action", "sub")
		form.Set("sr", "t5_test1,t5_test2,t5_test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.SubscribeByID(ctx, "t5_test1", "t5_test2", "t5_test3")
	require.NoError(t, err)
}

func TestSubredditService_Unsubscribe(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("action", "unsub")
		form.Set("sr_name", "test,golang,nba")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.Unsubscribe(ctx, "test", "golang", "nba")
	require.NoError(t, err)
}

func TestSubredditService_UnsubscribeByID(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("action", "unsub")
		form.Set("sr", "t5_test1,t5_test2,t5_test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.UnsubscribeByID(ctx, "t5_test1", "t5_test2", "t5_test3")
	require.NoError(t, err)
}

func TestSubredditService_Favorite(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/favorite", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("sr_name", "testsubreddit")
		form.Set("make_favorite", "true")
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.Favorite(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestSubredditService_Unfavorite(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/favorite", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("sr_name", "testsubreddit")
		form.Set("make_favorite", "false")
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.Unfavorite(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestSubredditService_Search(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/search", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("q", "golang")
		form.Set("limit", "10")
		form.Set("sort", "activity")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Search(ctx, "golang", &ListSubredditOptions{
		ListOptions: ListOptions{
			Limit: 10,
		},
		Sort: "activity",
	})
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_SearchNames(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/search-names.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/search_reddit_names", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("query", "golang")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	names, _, err := client.Subreddit.SearchNames(ctx, "golang")
	require.NoError(t, err)
	require.Equal(t, expectedSubredditNames, names)
}

func TestSubredditService_SearchPosts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/search-posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/all/search", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("q", "test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.SearchPosts(ctx, "test", "", nil)
	require.NoError(t, err)
	require.Equal(t, expectedSearchPosts, posts)
	require.Equal(t, "t3_hmwhd7", resp.After)
}

func TestSubredditService_SearchPosts_InSubreddit(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/search-posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/search", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("q", "test")
		form.Set("restrict_sr", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.SearchPosts(ctx, "test", "test", nil)
	require.NoError(t, err)
	require.Equal(t, expectedSearchPosts, posts)
	require.Equal(t, "t3_hmwhd7", resp.After)
}

func TestSubredditService_SearchPosts_InSubreddits(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/search-posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test+golang+nba/search", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("q", "test")
		form.Set("restrict_sr", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.SearchPosts(ctx, "test", "test+golang+nba", nil)
	require.NoError(t, err)
	require.Equal(t, expectedSearchPosts, posts)
	require.Equal(t, "t3_hmwhd7", resp.After)
}

func TestSubredditService_Random(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/random.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/random", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, "true", r.Form.Get("sr_detail"))
		require.Equal(t, "1", r.Form.Get("limit"))

		fmt.Fprint(w, blob)
	})

	subreddit, _, err := client.Subreddit.Random(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedRandomSubreddit, subreddit)
}

func TestSubredditService_RandomNSFW(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/random.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/randnsfw", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, "true", r.Form.Get("sr_detail"))
		require.Equal(t, "1", r.Form.Get("limit"))

		fmt.Fprint(w, blob)
	})

	subreddit, _, err := client.Subreddit.RandomNSFW(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedRandomSubreddit, subreddit)
}

func TestSubredditService_SubmissionText(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/test/api/submit_text", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `{
			"submit_text": "this is a test",
			"submit_text_html": ""
		}`)
	})

	text, _, err := client.Subreddit.SubmissionText(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, "this is a test", text)
}

func TestSubredditService_Banned(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/banned-users.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/banned", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("after", "testafter")
		form.Set("limit", "10")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	bans, _, err := client.Subreddit.Banned(ctx, "test", &ListOptions{After: "testafter", Limit: 10})
	require.NoError(t, err)
	require.Equal(t, expectedBans, bans)
}

func TestSubredditService_Muted(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/relationships.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/muted", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("before", "testbefore")
		form.Set("limit", "50")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	mutes, _, err := client.Subreddit.Muted(ctx, "test", &ListOptions{Before: "testbefore", Limit: 50})
	require.NoError(t, err)
	require.Equal(t, expectedRelationships3, mutes)
}

func TestSubredditService_WikiBanned(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/banned-users.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/wikibanned", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("after", "testafter")
		form.Set("limit", "15")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	bans, _, err := client.Subreddit.WikiBanned(ctx, "test", &ListOptions{After: "testafter", Limit: 15})
	require.NoError(t, err)
	require.Equal(t, expectedBans, bans)
}

func TestSubredditService_Contributors(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/relationships.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/contributors", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "5")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	contributors, _, err := client.Subreddit.Contributors(ctx, "test", &ListOptions{Limit: 5})
	require.NoError(t, err)
	require.Equal(t, expectedRelationships3, contributors)
}

func TestSubredditService_WikiContributors(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/relationships.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/wikicontributors", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "99")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	contributors, _, err := client.Subreddit.WikiContributors(ctx, "test", &ListOptions{Limit: 99})
	require.NoError(t, err)
	require.Equal(t, expectedRelationships3, contributors)
}

func TestSubredditService_Moderators(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/moderators.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/moderators", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	moderators, _, err := client.Subreddit.Moderators(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, expectedModerators, moderators)
}

func TestSubredditService_Rules(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/rules.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/rules", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	rules, _, err := client.Subreddit.Rules(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedRules, rules)
}

func TestSubredditService_CreateRule(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/add_subreddit_rule", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("kind", "all")
		form.Set("short_name", "testname")
		form.Set("violation_reason", "testreason")
		form.Set("description", "testdescription")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.CreateRule(ctx, "testsubreddit", &SubredditRuleCreateRequest{
		Kind:            "all",
		Name:            "testname",
		ViolationReason: "testreason",
		Description:     "testdescription",
	})
	require.NoError(t, err)
}

func TestSubredditService_CreateRule_Error(t *testing.T) {
	client, _ := setup(t)

	_, err := client.Subreddit.CreateRule(ctx, "testsubreddit", nil)
	require.EqualError(t, err, "*SubredditRuleCreateRequest: cannot be nil")

	_, err = client.Subreddit.CreateRule(ctx, "testsubreddit", &SubredditRuleCreateRequest{Kind: "invalid"})
	require.EqualError(t, err, "(*SubredditRuleCreateRequest).Kind: must be one of: comment, link, all")

	_, err = client.Subreddit.CreateRule(ctx, "testsubreddit", &SubredditRuleCreateRequest{Kind: "all", Name: ""})
	require.EqualError(t, err, "(*SubredditRuleCreateRequest).Name: must be between 1-100 characters")

	_, err = client.Subreddit.CreateRule(ctx, "testsubreddit", &SubredditRuleCreateRequest{
		Kind:            "all",
		Name:            "testname",
		ViolationReason: strings.Repeat("x", 101),
	})
	require.EqualError(t, err, "(*SubredditRuleCreateRequest).ViolationReason: cannot be longer than 100 characters")

	_, err = client.Subreddit.CreateRule(ctx, "testsubreddit", &SubredditRuleCreateRequest{
		Kind:        "all",
		Name:        "testname",
		Description: strings.Repeat("x", 501),
	})
	require.EqualError(t, err, "(*SubredditRuleCreateRequest).Description: cannot be longer than 500 characters")
}

func TestSubredditService_Traffic(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/traffic.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/traffic", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	dayTraffic, hourTraffic, monthTraffic, _, err := client.Subreddit.Traffic(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedDayTraffic, dayTraffic)
	require.Equal(t, expectedHourTraffic, hourTraffic)
	require.Equal(t, expectedMonthTraffic, monthTraffic)
}

func TestSubredditService_StyleSheet(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/stylesheet.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/stylesheet", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	styleSheet, _, err := client.Subreddit.StyleSheet(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedStyleSheet, styleSheet)
}

func TestSubredditService_StyleSheetRaw(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/stylesheet", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, "* { box-sizing: border-box; }")
	})

	styleSheet, _, err := client.Subreddit.StyleSheetRaw(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, "* { box-sizing: border-box; }", styleSheet)
}

func TestSubredditService_UpdateStyleSheet(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/subreddit_stylesheet", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("op", "save")
		form.Set("stylesheet_contents", "* { box-sizing: border-box; }")
		form.Set("reason", "testreason")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.UpdateStyleSheet(ctx, "testsubreddit", "* { box-sizing: border-box; }", "testreason")
	require.NoError(t, err)
}

func TestSubredditService_RemoveImage(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/delete_sr_img", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("img_name", "testimage")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.RemoveImage(ctx, "testsubreddit", "testimage")
	require.NoError(t, err)
}

func TestSubredditService_RemoveHeader(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/delete_sr_header", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.RemoveHeader(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestSubredditService_RemoveMobileHeader(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/delete_sr_banner", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.RemoveMobileHeader(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestSubredditService_RemoveMobileIcon(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/delete_sr_icon", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.RemoveMobileIcon(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestSubredditService_UploadImage(t *testing.T) {
	client, mux := setup(t)

	imageFile, err := ioutil.TempFile("/tmp", "emoji*.png")
	require.NoError(t, err)
	defer func() {
		imageFile.Close()
		os.Remove(imageFile.Name())
	}()

	_, err = imageFile.WriteString("this is a test")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/upload_sr_img", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		_, file, err := r.FormFile("file")
		require.NoError(t, err)

		rdr, err := file.Open()
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, rdr)
		require.NoError(t, err)
		require.Equal(t, "this is a test", buf.String())

		form := url.Values{}
		form.Set("upload_type", "img")
		form.Set("name", "testname")
		form.Set("img_type", "png")

		err = r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{
			"img_src": "https://example.com/test.png"
		}`)
	})

	link, _, err := client.Subreddit.UploadImage(ctx, "testsubreddit", imageFile.Name(), "testname")
	require.NoError(t, err)
	require.Equal(t, "https://example.com/test.png", link)
}

func TestSubredditService_UploadHeader(t *testing.T) {
	client, mux := setup(t)

	imageFile, err := ioutil.TempFile("/tmp", "emoji*.png")
	require.NoError(t, err)
	defer func() {
		imageFile.Close()
		os.Remove(imageFile.Name())
	}()

	_, err = imageFile.WriteString("this is a test")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/upload_sr_img", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		_, file, err := r.FormFile("file")
		require.NoError(t, err)

		rdr, err := file.Open()
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, rdr)
		require.NoError(t, err)
		require.Equal(t, "this is a test", buf.String())

		form := url.Values{}
		form.Set("upload_type", "header")
		form.Set("name", "testname")
		form.Set("img_type", "png")

		err = r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{
			"img_src": "https://example.com/test.png"
		}`)
	})

	link, _, err := client.Subreddit.UploadHeader(ctx, "testsubreddit", imageFile.Name(), "testname")
	require.NoError(t, err)
	require.Equal(t, "https://example.com/test.png", link)
}

func TestSubredditService_UploadMobileHeader(t *testing.T) {
	client, mux := setup(t)

	imageFile, err := ioutil.TempFile("/tmp", "emoji*.png")
	require.NoError(t, err)
	defer func() {
		imageFile.Close()
		os.Remove(imageFile.Name())
	}()

	_, err = imageFile.WriteString("this is a test")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/upload_sr_img", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		_, file, err := r.FormFile("file")
		require.NoError(t, err)

		rdr, err := file.Open()
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, rdr)
		require.NoError(t, err)
		require.Equal(t, "this is a test", buf.String())

		form := url.Values{}
		form.Set("upload_type", "banner")
		form.Set("name", "testname")
		form.Set("img_type", "png")

		err = r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{
			"img_src": "https://example.com/test.png"
		}`)
	})

	link, _, err := client.Subreddit.UploadMobileHeader(ctx, "testsubreddit", imageFile.Name(), "testname")
	require.NoError(t, err)
	require.Equal(t, "https://example.com/test.png", link)
}

func TestSubredditService_UploadMobileIcon(t *testing.T) {
	client, mux := setup(t)

	imageFile, err := ioutil.TempFile("/tmp", "emoji*.jpg")
	require.NoError(t, err)
	defer func() {
		imageFile.Close()
		os.Remove(imageFile.Name())
	}()

	_, err = imageFile.WriteString("this is a test")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/upload_sr_img", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		_, file, err := r.FormFile("file")
		require.NoError(t, err)

		rdr, err := file.Open()
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, rdr)
		require.NoError(t, err)
		require.Equal(t, "this is a test", buf.String())

		form := url.Values{}
		form.Set("upload_type", "icon")
		form.Set("name", "testname")
		form.Set("img_type", "jpg")

		err = r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{
			"img_src": "https://example.com/test.jpg"
		}`)
	})

	link, _, err := client.Subreddit.UploadMobileIcon(ctx, "testsubreddit", imageFile.Name(), "testname")
	require.NoError(t, err)
	require.Equal(t, "https://example.com/test.jpg", link)
}

func TestSubredditService_UploadImage_Error(t *testing.T) {
	client, mux := setup(t)

	imageFile, err := ioutil.TempFile("/tmp", "emoji*.jpg")
	require.NoError(t, err)
	defer func() {
		imageFile.Close()
		os.Remove(imageFile.Name())
	}()

	mux.HandleFunc("/r/testsubreddit/api/upload_sr_img", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		fmt.Fprint(w, `{
			"errors_values": [
				"error one",
				"error two"
			]
		}`)
	})

	_, _, err = client.Subreddit.UploadImage(ctx, "testsubreddit", "does-not-exist.jpg", "testname")
	require.EqualError(t, err, "open does-not-exist.jpg: no such file or directory")

	_, _, err = client.Subreddit.UploadImage(ctx, "testsubreddit", imageFile.Name(), "testname")
	require.EqualError(t, err, "could not upload image: error one; error two")
}

func TestSubredditService_Create(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/site_admin", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testsubreddit")
		form.Set("type", "private")
		form.Set("lang", "en")
		form.Set("title", "hello!")
		form.Set("public_description", "description")
		form.Set("description", "sidebar")
		form.Set("submit_text", "")
		form.Set("welcome_message_text", "")
		form.Set("welcome_message_enabled", "false")
		form.Set("allow_post_crossposts", "false")
		form.Set("allow_chat_post_creation", "true")
		form.Set("allow_polls", "false")
		form.Set("free_form_reports", "true")
		form.Set("original_content_tag_enabled", "false")
		form.Set("allow_images", "true")
		form.Set("allow_galleries", "true")
		form.Set("exclude_banned_modqueue", "false")
		form.Set("crowd_control_chat_level", "2")
		form.Set("all_original_content", "false")
		form.Set("submit_link_label", "submit a link!")
		form.Set("submit_text_label", "submit a post!")
		form.Set("link_type", "any")
		form.Set("spam_links", "low")
		form.Set("spam_selfposts", "low")
		form.Set("spam_comments", "low")
		form.Set("show_media", "false")
		form.Set("show_media_preview", "true")
		form.Set("collapse_deleted_comments", "false")
		form.Set("comment_score_hide_mins", "0")
		form.Set("spoilers_enabled", "true")
		form.Set("header-title", "hello!")
		form.Set("key_color", "")
		form.Set("hide_ads", "false")
		form.Set("over_18", "false")
		form.Set("allow_top", "true")
		form.Set("allow_discovery", "true")
		form.Set("wikimode", "modonly")
		form.Set("wiki_edit_age", "0")
		form.Set("wiki_edit_karma", "0")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.Create(ctx, "testsubreddit", nil)
	require.EqualError(t, err, "*SubredditSettings: cannot be nil")

	_, err = client.Subreddit.Create(ctx, "testsubreddit", expectedSubredditSettings)
	require.NoError(t, err)
}

func TestSubredditService_Edit(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/site_admin", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("sr", "t5_test")
		form.Set("type", "private")
		form.Set("lang", "en")
		form.Set("title", "hello!")
		form.Set("public_description", "description")
		form.Set("description", "sidebar")
		form.Set("submit_text", "")
		form.Set("welcome_message_text", "")
		form.Set("welcome_message_enabled", "false")
		form.Set("allow_post_crossposts", "false")
		form.Set("allow_chat_post_creation", "true")
		form.Set("allow_polls", "false")
		form.Set("free_form_reports", "true")
		form.Set("original_content_tag_enabled", "false")
		form.Set("allow_images", "true")
		form.Set("allow_galleries", "true")
		form.Set("exclude_banned_modqueue", "false")
		form.Set("crowd_control_chat_level", "2")
		form.Set("all_original_content", "false")
		form.Set("submit_link_label", "submit a link!")
		form.Set("submit_text_label", "submit a post!")
		form.Set("link_type", "any")
		form.Set("spam_links", "low")
		form.Set("spam_selfposts", "low")
		form.Set("spam_comments", "low")
		form.Set("show_media", "false")
		form.Set("show_media_preview", "true")
		form.Set("collapse_deleted_comments", "false")
		form.Set("comment_score_hide_mins", "0")
		form.Set("spoilers_enabled", "true")
		form.Set("header-title", "hello!")
		form.Set("key_color", "")
		form.Set("hide_ads", "false")
		form.Set("over_18", "false")
		form.Set("allow_top", "true")
		form.Set("allow_discovery", "true")
		form.Set("wikimode", "modonly")
		form.Set("wiki_edit_age", "0")
		form.Set("wiki_edit_karma", "0")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.Edit(ctx, "t5_test", nil)
	require.EqualError(t, err, "*SubredditSettings: cannot be nil")

	_, err = client.Subreddit.Edit(ctx, "t5_test", expectedSubredditSettings)
	require.NoError(t, err)
}

func TestSubredditService_GetSettings(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/settings.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/edit", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subredditSettings, _, err := client.Subreddit.GetSettings(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedSubredditSettings, subredditSettings)
}

func TestSubredditService_PostRequirements(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/post-requirements.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/v1/testsubreddit/post_requirements", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	postRequirements, _, err := client.Subreddit.PostRequirements(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedSubredditPostRequirements, postRequirements)
}
