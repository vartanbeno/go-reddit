package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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

func TestWikiService_Pages(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

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

	pages, _, err := client.Wiki.Pages(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, []string{"faq", "index"}, pages)
}

func TestWikiService_Settings(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	blob, err := readFileContents("../testdata/wiki/page-settings.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/wiki/settings/testpage", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	settings, _, err := client.Wiki.Settings(ctx, "testsubreddit", "testpage")
	require.NoError(t, err)
	require.Equal(t, expectedWikiPageSettings, settings)
}

func TestWikiService_UpdateSettings(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

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
	require.EqualError(t, err, "updateRequest: cannot be nil")

	settings, _, err := client.Wiki.UpdateSettings(ctx, "testsubreddit", "testpage", &WikiPageSettingsUpdateRequest{
		Listed:          Bool(false),
		PermissionLevel: PermissionApprovedContributorsOnly,
	})
	require.NoError(t, err)
	require.Equal(t, expectedWikiPageSettings, settings)
}

func TestWikiService_Allow(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

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
	client, mux, teardown := setup()
	defer teardown()

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
