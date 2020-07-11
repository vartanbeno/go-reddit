package reddit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/multi/multi.json")

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	multi, _, err := client.Multi.Get(ctx, "user/testuser/m/testmulti")
	assert.NoError(t, err)
	assert.Equal(t, expectedMulti, multi)
}

func TestMultiService_Mine(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/multi/multis.json")

	mux.HandleFunc("/api/multi/mine", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	multis, _, err := client.Multi.Mine(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []Multi{*expectedMulti, *expectedMulti2}, multis)
}

func TestMultiService_Of(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/multi/multis.json")

	mux.HandleFunc("/api/multi/user/test", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	multis, _, err := client.Multi.Of(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, []Multi{*expectedMulti, *expectedMulti2}, multis)
}

func TestMultiService_Copy(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/multi/multi.json")

	mux.HandleFunc("/api/multi/copy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("from", "user/testuser/m/testmulti")
		form.Set("to", "user/testuser2/m/testmulti2")
		form.Set("description_md", "this is a multireddit")
		form.Set("display_name", "hello")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	multi, _, err := client.Multi.Copy(ctx, &MultiCopyRequest{
		FromPath:    "user/testuser/m/testmulti",
		ToPath:      "user/testuser2/m/testmulti2",
		Description: "this is a multireddit",
		DisplayName: "hello",
	})
	assert.NoError(t, err)
	assert.Equal(t, expectedMulti, multi)
}

func TestMultiService_Create(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/multi/multi.json")
	createRequest := &MultiCreateOrUpdateRequest{
		Name:        "testmulti",
		Description: "this is a multireddit",
		Subreddits:  []string{"golang"},
		Visibility:  "public",
	}

	mux.HandleFunc("/api/multi", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		err := r.ParseForm()
		assert.NoError(t, err)

		model := r.Form.Get("model")

		expectedCreateRequest := new(MultiCreateOrUpdateRequest)
		err = json.Unmarshal([]byte(model), expectedCreateRequest)
		assert.NoError(t, err)
		assert.Equal(t, expectedCreateRequest, createRequest)

		fmt.Fprint(w, blob)
	})

	multi, _, err := client.Multi.Create(ctx, createRequest)
	assert.NoError(t, err)
	assert.Equal(t, expectedMulti, multi)
}

func TestMultiService_Update(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/multi/multi.json")
	updateRequest := &MultiCreateOrUpdateRequest{
		Name:        "testmulti",
		Description: "this is a multireddit",
		Visibility:  "public",
	}

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		err := r.ParseForm()
		assert.NoError(t, err)

		model := r.Form.Get("model")

		expectedCreateRequest := new(MultiCreateOrUpdateRequest)
		err = json.Unmarshal([]byte(model), expectedCreateRequest)
		assert.NoError(t, err)
		assert.Equal(t, expectedCreateRequest, updateRequest)

		fmt.Fprint(w, blob)
	})

	multi, _, err := client.Multi.Update(ctx, "user/testuser/m/testmulti", updateRequest)
	assert.NoError(t, err)
	assert.Equal(t, expectedMulti, multi)
}

func TestMultiService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
	})

	_, err := client.Multi.Delete(ctx, "user/testuser/m/testmulti")
	assert.NoError(t, err)
}

func TestMultiService_GetDescription(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/multi/description.json")

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti/description", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	description, _, err := client.Multi.GetDescription(ctx, "user/testuser/m/testmulti")
	assert.NoError(t, err)
	assert.Equal(t, "hello world", description)
}

func TestMultiService_UpdateDescription(t *testing.T) {
	setup()
	defer teardown()

	blob := readFileContents(t, "testdata/multi/description.json")

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti/description", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		form := url.Values{}
		form.Set("model", `{"body_md":"hello world"}`)

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	description, _, err := client.Multi.UpdateDescription(ctx, "user/testuser/m/testmulti", "hello world")
	assert.NoError(t, err)
	assert.Equal(t, "hello world", description)
}

func TestMultiService_AddSubreddit(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti/r/golang", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		form := url.Values{}
		form.Set("model", `{"name":"golang"}`)

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)
	})

	_, err := client.Multi.AddSubreddit(ctx, "user/testuser/m/testmulti", "golang")
	assert.NoError(t, err)
}

func TestMultiService_DeleteSubreddit(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/multi/user/testuser/m/testmulti/r/golang", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
	})

	_, err := client.Multi.DeleteSubreddit(ctx, "user/testuser/m/testmulti", "golang")
	assert.NoError(t, err)
}
