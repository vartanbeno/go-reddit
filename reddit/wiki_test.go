package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedWikiPage = &WikiPage{
	Content:   "test reason",
	Reason:    "this is a reason!",
	MayRevise: true,

	RevisionID:   "3c4e9fab-ef2c-11ea-90b6-0e9189256887",
	RevisionDate: &Timestamp{time.Date(2020, 9, 5, 3, 59, 45, 0, time.UTC)},
	RevisionBy: &User{
		ID:      "164ab8",
		Name:    "v_95",
		Created: &Timestamp{time.Date(2017, 3, 12, 4, 56, 47, 0, time.UTC)},

		PostKarma:    691,
		CommentKarma: 22235,

		HasVerifiedEmail: true,
		NSFW:             true,
	},
}

var expectedWikiPageSettings = &WikiPageSettings{
	PermissionLevel: PermissionSubredditWikiPermissions,
	Listed:          true,
	Editors: []*User{
		{
			ID:      "164ab8",
			Name:    "v_95",
			Created: &Timestamp{time.Date(2017, 3, 12, 4, 56, 47, 0, time.UTC)},

			PostKarma:    691,
			CommentKarma: 22235,

			HasVerifiedEmail: true,
			NSFW:             true,
		},
	},
}

var expectedWikiPageDiscussions = []*Post{
	{
		ID:      "imj8g5",
		FullID:  "t3_imj8g5",
		Created: &Timestamp{time.Date(2020, 9, 4, 16, 33, 33, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/helloworldtestt/comments/imj8g5/test/",
		URL:       "https://www.reddit.com/r/helloworldtestt/wiki/index",

		Title: "test",

		Likes: Bool(true),

		Score:            1,
		UpvoteRatio:      1,
		NumberOfComments: 0,

		SubredditName:         "helloworldtestt",
		SubredditNamePrefixed: "r/helloworldtestt",
		SubredditID:           "t5_2uquw1",
		SubredditSubscribers:  2,

		Author:   "v_95",
		AuthorID: "t2_164ab8",
	},
}

var expectedWikiPageRevisions = []*WikiPageRevision{
	{
		ID:      "3b28c343-effb-11ea-859e-0efe313b2cd3",
		Page:    "index",
		Created: &Timestamp{time.Date(2020, 9, 6, 4, 41, 29, 0, time.UTC)},
		Reason:  "reverted back 1 day",
		Hidden:  false,
		Author: &User{
			ID:      "164ab8",
			Name:    "v_95",
			Created: &Timestamp{time.Date(2017, 3, 12, 4, 56, 47, 0, time.UTC)},

			PostKarma:    691,
			CommentKarma: 22235,

			HasVerifiedEmail: true,
			NSFW:             true,
		},
	},
	{
		ID:      "d217592d-effa-11ea-90af-0e913d9ded0b",
		Page:    "index",
		Created: &Timestamp{time.Date(2020, 9, 6, 4, 38, 33, 0, time.UTC)},
		Reason:  "reverted back 10 minutes",
		Hidden:  false,
		Author: &User{
			ID:      "164ab8",
			Name:    "v_95",
			Created: &Timestamp{time.Date(2017, 3, 12, 4, 56, 47, 0, time.UTC)},

			PostKarma:    691,
			CommentKarma: 22235,

			HasVerifiedEmail: true,
			NSFW:             true,
		},
	},
}

func TestWikiService_Page(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/wiki/page.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/wiki/testpage", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	wikiPage, _, err := client.Wiki.Page(ctx, "testsubreddit", "testpage")
	require.NoError(t, err)
	require.Equal(t, expectedWikiPage, wikiPage)
}

func TestWikiService_PageRevision(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/wiki/page.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/wiki/testpage", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("v", "testrevision")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	wikiPage, _, err := client.Wiki.PageRevision(ctx, "testsubreddit", "testpage", "testrevision")
	require.NoError(t, err)
	require.Equal(t, expectedWikiPage, wikiPage)
}

func TestWikiService_Pages(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/wiki/pages", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `{
			"kind": "wikipagelisting",
			"data": [
				"faq",
				"index"
			]
		}`)
	})

	wikiPages, _, err := client.Wiki.Pages(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, []string{"faq", "index"}, wikiPages)
}

func TestWikiService_Edit(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/wiki/edit", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("page", "testpage")
		form.Set("content", "testcontent")
		form.Set("reason", "testreason")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Wiki.Edit(ctx, nil)
	require.EqualError(t, err, "*WikiPageEditRequest: cannot be nil")

	_, err = client.Wiki.Edit(ctx, &WikiPageEditRequest{
		Subreddit: "testsubreddit",
		Page:      "testpage",
		Content:   "testcontent",
		Reason:    "testreason",
	})
	require.NoError(t, err)
}

func TestWikiService_Revert(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/wiki/revert", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("page", "testpage")
		form.Set("revision", "testrevision")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Wiki.Revert(ctx, "testsubreddit", "testpage", "testrevision")
	require.NoError(t, err)
}

func TestWikiService_Settings(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/wiki/page-settings.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/wiki/settings/testpage", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	wikiPageSettings, _, err := client.Wiki.Settings(ctx, "testsubreddit", "testpage")
	require.NoError(t, err)
	require.Equal(t, expectedWikiPageSettings, wikiPageSettings)
}

func TestWikiService_UpdateSettings(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/wiki/page-settings.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/wiki/settings/testpage", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("permlevel", "1")
		form.Set("listed", "false")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	_, _, err = client.Wiki.UpdateSettings(ctx, "testsubreddit", "testpage", nil)
	require.EqualError(t, err, "*WikiPageSettingsUpdateRequest: cannot be nil")

	wikiPageSettings, _, err := client.Wiki.UpdateSettings(ctx, "testsubreddit", "testpage", &WikiPageSettingsUpdateRequest{
		Listed:          Bool(false),
		PermissionLevel: PermissionApprovedContributorsOnly,
	})
	require.NoError(t, err)
	require.Equal(t, expectedWikiPageSettings, wikiPageSettings)
}

func TestWikiService_Discussions(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/wiki/discussions.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/wiki/discussions/testpage", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	wikiPageDiscussions, _, err := client.Wiki.Discussions(ctx, "testsubreddit", "testpage", nil)
	require.NoError(t, err)
	require.Equal(t, expectedWikiPageDiscussions, wikiPageDiscussions)
}

func TestWikiService_ToggleVisibility(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/wiki/hide", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("page", "testpage")
		form.Set("revision", "testrevision")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{"status": true}`)
	})

	hidden, _, err := client.Wiki.ToggleVisibility(ctx, "testsubreddit", "testpage", "testrevision")
	require.NoError(t, err)
	require.True(t, hidden)
}

func TestWikiService_Revisions(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/wiki/revisions.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/wiki/revisions", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	wikiPageRevisions, _, err := client.Wiki.Revisions(ctx, "testsubreddit", nil)
	require.NoError(t, err)
	require.Equal(t, expectedWikiPageRevisions, wikiPageRevisions)
}

func TestWikiService_RevisionsPage(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/wiki/revisions.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/wiki/revisions/testpage", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "10")
		form.Set("after", "WikiRevision_wikiId1")
		form.Set("before", "WikiRevision_wikiId2")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	wikiPageRevisions, _, err := client.Wiki.RevisionsPage(ctx, "testsubreddit", "testpage", &ListOptions{
		Limit:  10,
		After:  "wikiId1",
		Before: "wikiId2",
	})
	require.NoError(t, err)
	require.Equal(t, expectedWikiPageRevisions, wikiPageRevisions)
}

func TestWikiService_Allow(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/wiki/alloweditor/add", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("page", "testpage")
		form.Set("username", "testusername")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Wiki.Allow(ctx, "testsubreddit", "testpage", "testusername")
	require.NoError(t, err)
}

func TestWikiService_Deny(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/wiki/alloweditor/del", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("page", "testpage")
		form.Set("username", "testusername")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Wiki.Deny(ctx, "testsubreddit", "testpage", "testusername")
	require.NoError(t, err)
}
