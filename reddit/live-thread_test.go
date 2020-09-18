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
	client, mux, teardown := setup()
	defer teardown()

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

func TestLiveThreadService_GetMultiple(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

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

func TestLiveThreadService_Create(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

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
	client, mux, teardown := setup()
	defer teardown()

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
	client, mux, teardown := setup()
	defer teardown()

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
	client, mux, teardown := setup()
	defer teardown()

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
	client, mux, teardown := setup()
	defer teardown()

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
	client, mux, teardown := setup()
	defer teardown()

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
	client, mux, teardown := setup()
	defer teardown()

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
	client, mux, teardown := setup()
	defer teardown()

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
	client, mux, teardown := setup()
	defer teardown()

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
	client, mux, teardown := setup()
	defer teardown()

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
	client, mux, teardown := setup()
	defer teardown()

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

	_, err := client.LiveThread.SetPermissions(ctx, "id123", "testuser", &LiveThreadPermissions{Discussions: true, Edit: true, Settings: true})
	require.NoError(t, err)
}

func TestLiveThreadService_Revoke(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

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

func TestLiveThreadService_Report(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

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
