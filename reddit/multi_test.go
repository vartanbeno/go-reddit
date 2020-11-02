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

var expectedMulti = &Multi{
	Name:        "test",
	DisplayName: "test",
	Path:        "/user/v_95/m/test/",
	Subreddits:  []string{"nba", "golang"},
	CopiedFrom:  nil,

	Owner:   "v_95",
	OwnerID: "t2_164ab8",
	Created: &Timestamp{time.Date(2020, 7, 11, 4, 55, 12, 0, time.UTC)},

	NumberOfSubscribers: 0,
	Visibility:          "private",
	CanEdit:             true,
}

var expectedMulti2 = &Multi{
	Name:        "test2",
	DisplayName: "test2",
	Path:        "/user/v_95/m/test2/",
	Subreddits:  []string{"redditdev", "test"},
	CopiedFrom:  nil,

	Owner:   "v_95",
	OwnerID: "t2_164ab8",
	Created: &Timestamp{time.Date(2020, 7, 11, 4, 57, 3, 0, time.UTC)},

	NumberOfSubscribers: 0,
	Visibility:          "private",
	CanEdit:             true,
}

func TestMultiService_Get(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/multi/multi.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	multi, _, err := client.Multi.Get(ctx, "user/testuser/m/testmulti")
	require.NoError(t, err)
	require.Equal(t, expectedMulti, multi)
}

func TestMultiService_Mine(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/multi/multis.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/multi/mine", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	multis, _, err := client.Multi.Mine(ctx)
	require.NoError(t, err)
	require.Equal(t, []*Multi{expectedMulti, expectedMulti2}, multis)
}

func TestMultiService_Of(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/multi/multis.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/multi/user/test", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	multis, _, err := client.Multi.Of(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, []*Multi{expectedMulti, expectedMulti2}, multis)
}

func TestMultiService_Copy(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/multi/multi.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/multi/copy", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("from", "user/testuser/m/testmulti")
		form.Set("to", "user/testuser2/m/testmulti2")
		form.Set("description_md", "this is a multireddit")
		form.Set("display_name", "hello")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	_, _, err = client.Multi.Copy(ctx, nil)
	require.EqualError(t, err, "*MultiCopyRequest: cannot be nil")

	multi, _, err := client.Multi.Copy(ctx, &MultiCopyRequest{
		FromPath:    "user/testuser/m/testmulti",
		ToPath:      "user/testuser2/m/testmulti2",
		Description: "this is a multireddit",
		DisplayName: "hello",
	})
	require.NoError(t, err)
	require.Equal(t, expectedMulti, multi)
}

func TestMultiService_Create(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/multi/multi.json")
	require.NoError(t, err)

	createRequest := &MultiCreateOrUpdateRequest{
		Name:        "testmulti",
		Description: "this is a multireddit",
		Subreddits:  []string{"golang"},
		Visibility:  "public",
	}

	mux.HandleFunc("/api/multi", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)

		model := r.Form.Get("model")

		expectedCreateRequest := new(MultiCreateOrUpdateRequest)
		err = json.Unmarshal([]byte(model), expectedCreateRequest)
		require.NoError(t, err)
		require.Equal(t, expectedCreateRequest, createRequest)

		fmt.Fprint(w, blob)
	})

	_, _, err = client.Multi.Create(ctx, nil)
	require.EqualError(t, err, "*MultiCreateOrUpdateRequest: cannot be nil")

	multi, _, err := client.Multi.Create(ctx, createRequest)
	require.NoError(t, err)
	require.Equal(t, expectedMulti, multi)
}

func TestMultiService_Update(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/multi/multi.json")
	require.NoError(t, err)

	updateRequest := &MultiCreateOrUpdateRequest{
		Name:        "testmulti",
		Description: "this is a multireddit",
		Visibility:  "public",
	}

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)

		model := r.Form.Get("model")

		expectedCreateRequest := new(MultiCreateOrUpdateRequest)
		err = json.Unmarshal([]byte(model), expectedCreateRequest)
		require.NoError(t, err)
		require.Equal(t, expectedCreateRequest, updateRequest)

		fmt.Fprint(w, blob)
	})

	_, _, err = client.Multi.Update(ctx, "user/testuser/m/testmulti", nil)
	require.EqualError(t, err, "*MultiCreateOrUpdateRequest: cannot be nil")

	multi, _, err := client.Multi.Update(ctx, "user/testuser/m/testmulti", updateRequest)
	require.NoError(t, err)
	require.Equal(t, expectedMulti, multi)
}

func TestMultiService_Delete(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
	})

	_, err := client.Multi.Delete(ctx, "user/testuser/m/testmulti")
	require.NoError(t, err)
}

func TestMultiService_Description(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/multi/description.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti/description", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	description, _, err := client.Multi.Description(ctx, "user/testuser/m/testmulti")
	require.NoError(t, err)
	require.Equal(t, "hello world", description)
}

func TestMultiService_UpdateDescription(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/multi/description.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti/description", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)

		form := url.Values{}
		form.Set("model", `{"body_md":"hello world"}`)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	description, _, err := client.Multi.UpdateDescription(ctx, "user/testuser/m/testmulti", "hello world")
	require.NoError(t, err)
	require.Equal(t, "hello world", description)
}

func TestMultiService_AddSubreddit(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti/r/golang", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)

		form := url.Values{}
		form.Set("model", `{"name":"golang"}`)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Multi.AddSubreddit(ctx, "user/testuser/m/testmulti", "golang")
	require.NoError(t, err)
}

func TestMultiService_DeleteSubreddit(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti/r/golang", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
	})

	_, err := client.Multi.DeleteSubreddit(ctx, "user/testuser/m/testmulti", "golang")
	require.NoError(t, err)
}
