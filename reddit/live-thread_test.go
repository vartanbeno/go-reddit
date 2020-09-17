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
	require.EqualError(t, err, "*LiveThreadCreateRequest: cannot be nil")

	id, _, err := client.LiveThread.Create(ctx, &LiveThreadCreateRequest{
		Title:       "testtitle",
		Description: "testdescription",
		Resources:   "testresources",
		NSFW:        true,
	})
	require.NoError(t, err)
	require.Equal(t, "id123", id)
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
