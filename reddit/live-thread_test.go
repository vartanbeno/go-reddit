package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedLiveThread = &LiveThread{
	ID:      "15nevtv8e54dh",
	FullID:  "LiveUpdateEvent_15nevtv8e54dh",
	Created: &Timestamp{time.Date(2020, 9, 16, 1, 20, 27, 0, time.UTC)},

	Title:       "test",
	Description: "test",
	Resources:   "",

	State:             "live",
	ViewerCount:       6,
	ViewerCountFuzzed: true,

	WebSocketURL: "wss://ws-078adc7cb2099a9df.wss.redditmedia.com/live/15nevtv8e54dh?m=AQAA7rxiX6EpLYFCFZ0KJD4lVAPaMt0A1z2-xJ1b2dWCmxNIfMwL",

	Announcement: false,
	NSFW:         false,
}

var expectedLiveThreads = []*LiveThread{
	{
		ID:      "15nevtv8e54dh",
		FullID:  "LiveUpdateEvent_15nevtv8e54dh",
		Created: &Timestamp{time.Date(2020, 9, 16, 1, 20, 27, 0, time.UTC)},

		Title:       "test",
		Description: "test",
		Resources:   "",

		State:             "live",
		ViewerCount:       6,
		ViewerCountFuzzed: true,

		WebSocketURL: "wss://ws-078adc7cb2099a9df.wss.redditmedia.com/live/15nevtv8e54dh?m=AQAA7rxiX6EpLYFCFZ0KJD4lVAPaMt0A1z2-xJ1b2dWCmxNIfMwL",

		Announcement: false,
		NSFW:         false,
	},
	{
		ID:      "15ndkho8e54dh",
		FullID:  "LiveUpdateEvent_15ndkho8e54dh",
		Created: &Timestamp{time.Date(2020, 9, 16, 1, 20, 37, 0, time.UTC)},

		Title:       "test 2",
		Description: "test 2",
		Resources:   "",

		State:             "live",
		ViewerCount:       6,
		ViewerCountFuzzed: true,

		WebSocketURL: "wss://ws-078adc7cb2099a9df.wss.redditmedia.com/live/15ndkho8e54dh?m=AQAA7rxiX6EpLYFCFZ0KJD4lVAPaMt0A1z2-xJ1b2dWCmxNIfMwL",

		Announcement: false,
		NSFW:         false,
	},
}

var expectedLiveThreadUpdates = []*LiveThreadUpdate{
	{
		ID:      "5e46cd94-f968-11ea-9a6a-0e1933241e7d",
		FullID:  "LiveUpdate_5e46cd94-f968-11ea-9a6a-0e1933241e7d",
		Author:  "testuser1",
		Created: &Timestamp{time.Date(2020, 9, 18, 4, 35, 24, 0, time.UTC)},

		Body:         "test 2",
		EmbeddedURLs: []string{"https://example.com", "https://reddit.com"},

		Stricken: true,
	},
	{
		ID:      "fc44f204-f964-11ea-b148-0e2e56a0425f",
		FullID:  "LiveUpdate_fc44f204-f964-11ea-b148-0e2e56a0425f",
		Author:  "testuser1",
		Created: &Timestamp{time.Date(2020, 9, 18, 4, 11, 11, 0, time.UTC)},

		Body: "test 1",

		Stricken: true,
	},
}

var expectedLiveThreadUpdate = &LiveThreadUpdate{
	ID:      "fc44f204-f964-11ea-b148-0e2e56a0425f",
	FullID:  "LiveUpdate_fc44f204-f964-11ea-b148-0e2e56a0425f",
	Author:  "testuser1",
	Created: &Timestamp{time.Date(2020, 9, 18, 4, 11, 11, 0, time.UTC)},

	Body: "test 1",

	Stricken: true,
}

var expectedLiveThreadDiscussions = []*Post{
	{
		ID:      "test1",
		FullID:  "t3_test1",
		Created: &Timestamp{time.Date(2020, 9, 16, 12, 37, 31, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/live/comments/test1/test_title/",
		URL:       "https://www.reddit.com/live/15nfp4mtfbo14/",

		Title: "test title",

		Score:            22,
		UpvoteRatio:      0.9,
		NumberOfComments: 1,

		SubredditName:         "live",
		SubredditNamePrefixed: "r/live",
		SubredditID:           "t5_32o7w",
		SubredditSubscribers:  24501,

		Author:   "TestUser",
		AuthorID: "t2_test1",
	},
	{
		ID:      "test2",
		FullID:  "t3_test2",
		Created: &Timestamp{time.Date(2020, 9, 16, 12, 37, 1, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/live/comments/test2/test_title/",
		URL:       "https://www.reddit.com/live/15nfp4mtfbo14/",

		Title: "test title",

		Score:            71,
		UpvoteRatio:      0.97,
		NumberOfComments: 34,

		SubredditName:         "live",
		SubredditNamePrefixed: "r/live",
		SubredditID:           "t5_32o7w",
		SubredditSubscribers:  24501,

		Author:   "TestUser",
		AuthorID: "t2_test1",
	},
}

var expectedLiveThreadContributors = &LiveThreadContributors{
	Current: []*LiveThreadContributor{
		{ID: "t2_test1", Name: "test1", Permissions: []string{"all"}},
		{ID: "t2_test2", Name: "test2", Permissions: []string{"all"}},
	},
	Invited: nil,
}

var expectedLiveThreadContributorsAndInvited = &LiveThreadContributors{
	Current: []*LiveThreadContributor{
		{ID: "t2_test1", Name: "test1", Permissions: []string{"all"}},
		{ID: "t2_test2", Name: "test2", Permissions: []string{"all"}},
	},
	Invited: []*LiveThreadContributor{
		{ID: "t2_test3", Name: "test3", Permissions: []string{"manage", "discussions"}},
	},
}

func TestLiveThreadService_Get(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/live-thread/live-thread.json")
	require.NoError(t, err)

	mux.HandleFunc("/live/id123/about", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	liveThread, _, err := client.LiveThread.Get(ctx, "id123")
	require.NoError(t, err)
	require.Equal(t, expectedLiveThread, liveThread)
}

func TestLiveThreadService_Now(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/live-thread/live-thread.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/live/happening_now", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	liveThread, _, err := client.LiveThread.Now(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedLiveThread, liveThread)
}

func TestLiveThreadService_Now_NoContent(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/happening_now", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(204)
	})

	liveThread, _, err := client.LiveThread.Now(ctx)
	require.NoError(t, err)
	require.Nil(t, liveThread)
}

func TestLiveThreadService_GetMultiple(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/live-thread/live-threads.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/live/by_id/id1,id2", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	_, _, err = client.LiveThread.GetMultiple(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	liveThreads, _, err := client.LiveThread.GetMultiple(ctx, "id1", "id2")
	require.NoError(t, err)
	require.Equal(t, expectedLiveThreads, liveThreads)
}

func TestLiveThreadService_Update(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/update", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("body", "test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.Update(ctx, "id123", "test")
	require.NoError(t, err)
}

func TestLiveThreadService_Updates(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/live-thread/updates.json")
	require.NoError(t, err)

	mux.HandleFunc("/live/id123", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	liveThreadUpdates, _, err := client.LiveThread.Updates(ctx, "id123", nil)
	require.NoError(t, err)
	require.Equal(t, expectedLiveThreadUpdates, liveThreadUpdates)
}

func TestLiveThreadService_UpdateByID(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/live-thread/update.json")
	require.NoError(t, err)

	mux.HandleFunc("/live/id123/updates/update123", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	liveThreadUpdate, _, err := client.LiveThread.UpdateByID(ctx, "id123", "update123")
	require.NoError(t, err)
	require.Equal(t, expectedLiveThreadUpdate, liveThreadUpdate)
}

func TestLiveThreadService_Discussions(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/live-thread/discussions.json")
	require.NoError(t, err)

	mux.HandleFunc("/live/id123/discussions", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	liveThreadDiscussions, _, err := client.LiveThread.Discussions(ctx, "id123", nil)
	require.NoError(t, err)
	require.Equal(t, expectedLiveThreadDiscussions, liveThreadDiscussions)
}

func TestLiveThreadService_Strike(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/strike_update", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "update123")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.Strike(ctx, "id123", "update123")
	require.NoError(t, err)
}

func TestLiveThreadService_Delete(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/delete_update", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "update123")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.Delete(ctx, "id123", "update123")
	require.NoError(t, err)
}

func TestLiveThreadService_Create(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/create", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("title", "testtitle")
		form.Set("description", "testdescription")
		form.Set("resources", "testresources")
		form.Set("nsfw", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{
			"json": {
				"data": {
					"id": "id123"
				},
				"errors": []
			}
		}`)
	})

	_, _, err := client.LiveThread.Create(ctx, nil)
	require.EqualError(t, err, "*LiveThreadCreateOrUpdateRequest: cannot be nil")

	id, _, err := client.LiveThread.Create(ctx, &LiveThreadCreateOrUpdateRequest{
		Title:       "testtitle",
		Description: "testdescription",
		Resources:   "testresources",
		NSFW:        Bool(true),
	})
	require.NoError(t, err)
	require.Equal(t, "id123", id)
}

func TestLiveThreadService_Close(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/close_thread", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.Close(ctx, "id123")
	require.NoError(t, err)
}

func TestLiveThreadService_Configure(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/edit", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("title", "testtitle")
		form.Set("description", "testdescription")
		form.Set("resources", "testresources")
		form.Set("nsfw", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{
			"json": {
				"data": {
					"id": "id123"
				},
				"errors": []
			}
		}`)
	})

	_, err := client.LiveThread.Configure(ctx, "id123", nil)
	require.EqualError(t, err, "*LiveThreadCreateOrUpdateRequest: cannot be nil")

	_, err = client.LiveThread.Configure(ctx, "id123", &LiveThreadCreateOrUpdateRequest{
		Title:       "testtitle",
		Description: "testdescription",
		Resources:   "testresources",
		NSFW:        Bool(true),
	})
	require.NoError(t, err)
}

func TestLiveThreadService_Contributors(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/live-thread/contributors.json")
	require.NoError(t, err)

	mux.HandleFunc("/live/id123/contributors", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	contributors, _, err := client.LiveThread.Contributors(ctx, "id123")
	require.NoError(t, err)
	require.Equal(t, expectedLiveThreadContributors, contributors)
}

func TestLiveThreadService_ContributorsAndInvited(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/live-thread/contributors-and-invited.json")
	require.NoError(t, err)

	mux.HandleFunc("/live/id123/contributors", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	contributors, _, err := client.LiveThread.Contributors(ctx, "id123")
	require.NoError(t, err)
	require.Equal(t, expectedLiveThreadContributorsAndInvited, contributors)
}

func TestLiveThreadService_Accept(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/accept_contributor_invite", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.Accept(ctx, "id123")
	require.NoError(t, err)
}

func TestLiveThreadService_Leave(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/leave_contributor", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.Leave(ctx, "id123")
	require.NoError(t, err)
}

func TestLiveThreadService_Invite(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/invite_contributor", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "liveupdate_contributor_invite")
		form.Set("permissions", "+all")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.Invite(ctx, "id123", "testuser", nil)
	require.NoError(t, err)
}

func TestLiveThreadService_Invite_Permissions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/invite_contributor", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "liveupdate_contributor_invite")
		form.Set("permissions", "-all,+close,-discussions,-edit,+manage,-settings,+update")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.Invite(ctx, "id123", "testuser", &LiveThreadPermissions{Close: true, Manage: true, Update: true})
	require.NoError(t, err)
}

func TestLiveThreadService_Uninvite(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/rm_contributor_invite", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t2_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.Uninvite(ctx, "id123", "t2_test")
	require.NoError(t, err)
}

func TestLiveThreadService_SetPermissions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/set_contributor_permissions", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "liveupdate_contributor")
		form.Set("permissions", "-all,-close,+discussions,+edit,-manage,+settings,-update")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.SetPermissions(ctx, "id123", "testuser", &LiveThreadPermissions{Discussions: true, Edit: true, Settings: true})
	require.NoError(t, err)
}

func TestLiveThreadService_SetPermissionsForInvite(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/set_contributor_permissions", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "liveupdate_contributor_invite")
		form.Set("permissions", "-all,-close,+discussions,+edit,-manage,+settings,-update")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.SetPermissionsForInvite(ctx, "id123", "testuser", &LiveThreadPermissions{Discussions: true, Edit: true, Settings: true})
	require.NoError(t, err)
}

func TestLiveThreadService_Revoke(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/rm_contributor", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("id", "t2_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.Revoke(ctx, "id123", "t2_test")
	require.NoError(t, err)
}

func TestLiveThreadService_HideDiscussion(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/hide_discussion", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("link", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.HideDiscussion(ctx, "id123", "t3_test")
	require.NoError(t, err)
}

func TestLiveThreadService_UnhideDiscussion(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/unhide_discussion", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("link", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.UnhideDiscussion(ctx, "id123", "t3_test")
	require.NoError(t, err)
}

func TestLiveThreadService_Report(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/live/id123/report", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("type", "spam")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.LiveThread.Report(ctx, "id123", "invalidreason")
	require.EqualError(t, err, "invalid reason for reporting live thread: invalidreason")

	_, err = client.LiveThread.Report(ctx, "id123", "spam")
	require.NoError(t, err)
}
