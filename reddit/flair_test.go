package reddit

import (
	"encoding/json"
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

var expectedFlairTemplate = &FlairTemplate{
	ID:      "be0a6110-f23c-11ea-862f-0e08890d7323",
	Type:    "LINK_FLAIR",
	ModOnly: false,

	AllowableContent: "all",
	Text:             "lol",
	TextType:         "richtext",
	TextColor:        "dark",
	TextEditable:     false,
	RichText: []map[string]string{
		{"e": "text", "t": "lol"},
	},

	OverrideCSS:     false,
	MaxEmojis:       1,
	BackgroundColor: "#fafafa",
	CSSClass:        "",
}

var expectedFlairChoices = []*FlairChoice{
	{
		TemplateID: "c4edd5ce-40e8-11e7-b814-0ef91bd65558",
		Text:       "Reddit API",
		Editable:   false,
		Position:   "left",
		CSSClass:   "",
	},
	{
		TemplateID: "49bb3d06-0dad-11e7-b897-0e42c2400b7a",
		Text:       "PRAW",
		Editable:   false,
		Position:   "left",
		CSSClass:   "",
	},
	{
		TemplateID: "f1905376-40e9-11e7-a0dc-0e2f53ef3712",
		Text:       "snoowrap",
		Editable:   false,
		Position:   "left",
		CSSClass:   "",
	},
	{
		TemplateID: "03dc6ea8-40e9-11e7-8abb-0eb85aed0bce",
		Text:       "Other API Wrapper",
		Editable:   false,
		Position:   "left",
		CSSClass:   "",
	},
}

var expectedFlairChoice = &FlairChoice{
	TemplateID: "03dc6ea8-40e9-11e7-8abb-0eb85aed0bce",
	Text:       "Other API Wrapper",
	Editable:   false,
	Position:   "left",
	CSSClass:   "",
}

var expectedFlairChanges = []*FlairChangeResponse{
	{
		OK:       false,
		Status:   "skipped",
		Warnings: map[string]string{},
		Errors: map[string]string{
			"user": "unable to resolve user `testuser1', ignoring",
		},
	},
	{
		OK:       true,
		Status:   "added flair for user testuser2",
		Warnings: map[string]string{},
		Errors:   map[string]string{},
	},
	{
		OK:       true,
		Status:   "added flair for user testuser3",
		Warnings: map[string]string{},
		Errors:   map[string]string{},
	},
	{
		OK:       true,
		Status:   "removed flair for user testuser4",
		Warnings: map[string]string{},
		Errors:   map[string]string{},
	},
}

func TestFlairService_GetUserFlairs(t *testing.T) {
	client, mux := setup(t)

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
	client, mux := setup(t)

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
	client, mux := setup(t)

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
	client, mux := setup(t)

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
	require.EqualError(t, err, "*FlairConfigureRequest: cannot be nil")

	_, err = client.Flair.Configure(ctx, "testsubreddit", &FlairConfigureRequest{
		UserFlairEnabled:           Bool(true),
		UserFlairPosition:          "right",
		UserFlairSelfAssignEnabled: Bool(false),
		PostFlairPosition:          "left",
		PostFlairSelfAssignEnabled: Bool(false),
	})
	require.NoError(t, err)
}

func TestFlairService_Enable(t *testing.T) {
	client, mux := setup(t)

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
	client, mux := setup(t)

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

func TestFlairService_UpsertUserTemplate(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/flair/flair-template.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/flairtemplate_v2", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("flair_type", "USER_FLAIR")
		form.Set("allowable_content", "all")
		form.Set("text", "testtext")
		form.Set("text_color", "dark")
		form.Set("text_editable", "false")
		form.Set("mod_only", "true")
		form.Set("max_emojis", "5")
		form.Set("background_color", "transparent")
		form.Set("css_class", "testclass")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	_, _, err = client.Flair.UpsertUserTemplate(ctx, "testsubreddit", nil)
	require.EqualError(t, err, "*FlairTemplateCreateOrUpdateRequest: cannot be nil")

	flairTemplate, _, err := client.Flair.UpsertUserTemplate(ctx, "testsubreddit", &FlairTemplateCreateOrUpdateRequest{
		AllowableContent: "all",
		ModOnly:          Bool(true),
		Text:             "testtext",
		TextColor:        "dark",
		TextEditable:     Bool(false),
		MaxEmojis:        Int(5),
		BackgroundColor:  "transparent",
		CSSClass:         "testclass",
	})
	require.NoError(t, err)
	require.Equal(t, expectedFlairTemplate, flairTemplate)
}

func TestFlairService_UpsertPostTemplate(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/flair/flair-template.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/flairtemplate_v2", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("flair_type", "LINK_FLAIR")
		form.Set("flair_template_id", "testid")
		form.Set("allowable_content", "text")
		form.Set("text", "testtext")
		form.Set("text_color", "light")
		form.Set("text_editable", "true")
		form.Set("mod_only", "false")
		form.Set("background_color", "#fafafa")
		form.Set("css_class", "testclass")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	_, _, err = client.Flair.UpsertPostTemplate(ctx, "testsubreddit", nil)
	require.EqualError(t, err, "*FlairTemplateCreateOrUpdateRequest: cannot be nil")

	flairTemplate, _, err := client.Flair.UpsertPostTemplate(ctx, "testsubreddit", &FlairTemplateCreateOrUpdateRequest{
		ID:               "testid",
		AllowableContent: "text",
		ModOnly:          Bool(false),
		Text:             "testtext",
		TextColor:        "light",
		TextEditable:     Bool(true),
		BackgroundColor:  "#fafafa",
		CSSClass:         "testclass",
	})
	require.NoError(t, err)
	require.Equal(t, expectedFlairTemplate, flairTemplate)
}

func TestFlairService_Delete(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/deleteflair", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Flair.Delete(ctx, "testsubreddit", "testuser")
	require.NoError(t, err)
}

func TestFlairService_DeleteTemplate(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/deleteflairtemplate", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("flair_template_id", "testtemplate")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Flair.DeleteTemplate(ctx, "testsubreddit", "testtemplate")
	require.NoError(t, err)
}

func TestFlairService_DeleteAllUserTemplates(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/clearflairtemplates", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("flair_type", "USER_FLAIR")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Flair.DeleteAllUserTemplates(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestFlairService_DeleteAllPostTemplates(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/clearflairtemplates", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("flair_type", "LINK_FLAIR")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Flair.DeleteAllPostTemplates(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestFlairService_ReorderUserTemplates(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/testsubreddit/flair_template_order/USER_FLAIR", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)

		var ids []string
		err := json.NewDecoder(r.Body).Decode(&ids)
		require.NoError(t, err)
		require.Equal(t, []string{"test1", "test2", "test3", "test4"}, ids)
	})

	_, err := client.Flair.ReorderUserTemplates(ctx, "testsubreddit", []string{"test1", "test2", "test3", "test4"})
	require.NoError(t, err)
}

func TestFlairService_ReorderPostTemplates(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/testsubreddit/flair_template_order/LINK_FLAIR", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)

		var ids []string
		err := json.NewDecoder(r.Body).Decode(&ids)
		require.NoError(t, err)
		require.Equal(t, []string{"test1", "test2", "test3", "test4"}, ids)
	})

	_, err := client.Flair.ReorderPostTemplates(ctx, "testsubreddit", []string{"test1", "test2", "test3", "test4"})
	require.NoError(t, err)
}

func TestFlairService_Choices(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/flair/choices.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/flairselector", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("name", "user1")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	choices, current, _, err := client.Flair.Choices(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedFlairChoices, choices)
	require.Equal(t, expectedFlairChoice, current)
}

func TestFlairService_ChoicesOf(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/flair/choices.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/flairselector", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("name", "testuser")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	choices, current, _, err := client.Flair.ChoicesOf(ctx, "testsubreddit", "testuser")
	require.NoError(t, err)
	require.Equal(t, expectedFlairChoices, choices)
	require.Equal(t, expectedFlairChoice, current)
}

func TestFlairService_ChoicesForPost(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/flair/choices.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/flairselector", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("link", "t3_123")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	choices, current, _, err := client.Flair.ChoicesForPost(ctx, "t3_123")
	require.NoError(t, err)
	require.Equal(t, expectedFlairChoices, choices)
	require.Equal(t, expectedFlairChoice, current)
}

func TestFlairService_ChoicesForNewPost(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/flair/choices.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/flairselector", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("is_newlink", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	choices, _, err := client.Flair.ChoicesForNewPost(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedFlairChoices, choices)
}

func TestFlairService_Select(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/selectflair", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "user1")
		form.Set("flair_template_id", "id123")
		form.Set("text", "text123")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Flair.Select(ctx, "testsubreddit", nil)
	require.EqualError(t, err, "*FlairSelectRequest: cannot be nil")

	_, err = client.Flair.Select(ctx, "testsubreddit", &FlairSelectRequest{
		ID:   "id123",
		Text: "text123",
	})
	require.NoError(t, err)
}

func TestFlairService_Assign(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/selectflair", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("flair_template_id", "id123")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Flair.Assign(ctx, "testsubreddit", "testuser", nil)
	require.EqualError(t, err, "*FlairSelectRequest: cannot be nil")

	_, err = client.Flair.Assign(ctx, "testsubreddit", "testuser", &FlairSelectRequest{
		ID: "id123",
	})
	require.NoError(t, err)
}

func TestFlairService_SelectForPost(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/selectflair", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("link", "t3_123")
		form.Set("flair_template_id", "id123")
		form.Set("text", "text123")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Flair.SelectForPost(ctx, "t3_123", nil)
	require.EqualError(t, err, "*FlairSelectRequest: cannot be nil")

	_, err = client.Flair.SelectForPost(ctx, "t3_123", &FlairSelectRequest{
		ID:   "id123",
		Text: "text123",
	})
	require.NoError(t, err)
}

func TestFlairService_RemoveFromPost(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/selectflair", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("link", "t3_123")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Flair.RemoveFromPost(ctx, "t3_123")
	require.NoError(t, err)
}

func TestFlairService_Change(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/flair/csv-change.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/flaircsv", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("flair_csv", `testuser1,testtext1,testclass1
testuser2,testtext2,testclass2
testuser3,testtext3,testclass3
testuser4,testtext4,testclass4
`)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	_, _, err = client.Flair.Change(ctx, "testsubreddit", nil)
	require.EqualError(t, err, "requests: must provide between 1 and 100")

	changes, _, err := client.Flair.Change(ctx, "testsubreddit", []FlairChangeRequest{
		{"testuser1", "testtext1", "testclass1"},
		{"testuser2", "testtext2", "testclass2"},
		{"testuser3", "testtext3", "testclass3"},
		{"testuser4", "testtext4", "testclass4"},
	})
	require.NoError(t, err)
	require.Equal(t, expectedFlairChanges, changes)
}
