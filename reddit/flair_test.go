package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

var expectedUserFlairs = []*Flair{
	{
		ID:   "b8a1c822-3feb-11e8-88e1-0e5f55d58ce0",
		Type: "text",
		Text: "Beginner",

		Color:           "dark",
		BackgroundColor: "",
		CSSClass:        "Beginner1",

		Editable: false,
		ModOnly:  false,
	},
	{
		ID:   "b8ea0fce-3feb-11e8-af7a-0e263a127cf8",
		Text: "Moderator",
		Type: "text",

		Color:           "dark",
		BackgroundColor: "",
		CSSClass:        "Moderator1",

		Editable: false,
		ModOnly:  true,
	},
}

var expectedPostFlairs = []*Flair{
	{
		ID:   "305b503e-da60-11ea-9681-0e9f1d580d2d",
		Type: "richtext",
		Text: "test",

		Color:           "light",
		BackgroundColor: "#373c3f",
		CSSClass:        "test",

		Editable: false,
		ModOnly:  true,
	},
}

var expectedListUserFlairs = []*FlairSummary{
	{
		User: "TestUser1",
		Text: "TestFlair1",
	},
	{
		User: "TestUser2",
		Text: "TestFlair2",
	},
}

func TestFlairService_GetUserFlairs(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	blob, err := readFileContents("../testdata/flair/user-flairs.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/user_flair_v2", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	userFlairs, _, err := client.Flair.GetUserFlairs(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedUserFlairs, userFlairs)
}

func TestFlairService_GetPostFlairs(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	blob, err := readFileContents("../testdata/flair/post-flairs.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/link_flair_v2", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	postFlairs, _, err := client.Flair.GetPostFlairs(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedPostFlairs, postFlairs)
}

func TestFlairService_ListUserFlairs(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	blob, err := readFileContents("../testdata/flair/list-user-flairs.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/flairlist", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	userFlairs, _, err := client.Flair.ListUserFlairs(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedListUserFlairs, userFlairs)
}

func TestFlairService_Configure(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/r/testsubreddit/api/flairconfig", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("flair_enabled", "true")
		form.Set("flair_position", "right")
		form.Set("flair_self_assign_enabled", "false")
		form.Set("link_flair_position", "left")
		form.Set("link_flair_self_assign_enabled", "false")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Flair.Configure(ctx, "testsubreddit", nil)
	require.EqualError(t, err, "request: cannot be nil")

	_, err = client.Flair.Configure(ctx, "testsubreddit", &ConfigureFlairRequest{
		UserFlairEnabled:           Bool(true),
		UserFlairPosition:          "right",
		UserFlairSelfAssignEnabled: Bool(false),
		PostFlairPosition:          "left",
		PostFlairSelfAssignEnabled: Bool(false),
	})
	require.NoError(t, err)
}

func TestFlairService_Enable(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/r/testsubreddit/api/setflairenabled", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("flair_enabled", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Flair.Enable(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestFlairService_Disable(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/r/testsubreddit/api/setflairenabled", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("flair_enabled", "false")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Flair.Disable(ctx, "testsubreddit")
	require.NoError(t, err)
}
