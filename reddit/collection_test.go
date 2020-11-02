package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedCollection = &Collection{
	ID:      "37f1e52d-7ec9-466b-b4cc-59e86e071ed7",
	Created: &Timestamp{time.Date(2020, 8, 6, 23, 25, 3, 0, time.UTC)},
	Updated: &Timestamp{time.Date(2020, 8, 7, 1, 59, 32, 0, time.UTC)},

	Title:     "Test Title",
	Permalink: "https://www.reddit.com/r/helloworldtestt/collection/37f1e52d-7ec9-466b-b4cc-59e86e071ed7",
	Layout:    "TIMELINE",

	SubredditID:   "t5_2uquw1",
	Author:        "v_95",
	AuthorID:      "t2_164ab8",
	PrimaryPostID: "t3_hs0cyh",
	PostIDs:       []string{"t3_hs0cyh", "t3_hqrg8s", "t3_hs03f3"},
}

var expectedCollections = []*Collection{
	{
		ID:      "37f1e52d-7ec9-466b-b4cc-59e86e071ed7",
		Created: &Timestamp{time.Date(2020, 8, 6, 23, 25, 3, 0, time.UTC)},
		Updated: &Timestamp{time.Date(2020, 8, 7, 1, 59, 32, 0, time.UTC)},

		Title:     "Test Title",
		Permalink: "https://www.reddit.com/r/helloworldtestt/collection/37f1e52d-7ec9-466b-b4cc-59e86e071ed7",
		Layout:    "TIMELINE",

		SubredditID: "t5_2uquw1",
		Author:      "v_95",
		AuthorID:    "t2_164ab8",
		PostIDs:     []string{"t3_hs0cyh", "t3_hqrg8s", "t3_hs03f3"},
	},
	{
		ID:      "8e94db00-6605-46c6-b0d2-44653d6f538c",
		Created: &Timestamp{time.Date(2020, 8, 7, 0, 56, 29, 0, time.UTC)},
		Updated: &Timestamp{time.Date(2020, 8, 7, 1, 59, 27, 0, time.UTC)},

		Title:       "Test Title 2",
		Description: "Test Description",
		Permalink:   "https://www.reddit.com/r/helloworldtestt/collection/8e94db00-6605-46c6-b0d2-44653d6f538c",

		SubredditID: "t5_2uquw1",
		Author:      "v_95",
		AuthorID:    "t2_164ab8",
		PostIDs:     []string{},
	},
	{
		ID:      "a1b3e088-f6b8-4d98-9e93-adaacef113cd",
		Created: &Timestamp{time.Date(2020, 8, 7, 0, 55, 24, 0, time.UTC)},
		Updated: &Timestamp{time.Date(2020, 8, 7, 0, 55, 24, 0, time.UTC)},

		Title:     "Test Title 3",
		Permalink: "https://www.reddit.com/r/helloworldtestt/collection/a1b3e088-f6b8-4d98-9e93-adaacef113cd",

		SubredditID: "t5_2uquw1",
		Author:      "v_95",
		AuthorID:    "t2_164ab8",
		PostIDs:     []string{},
	},
}

func TestCollectionService_Get(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/collection/collection.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/v1/collections/collection", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("collection_id", "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
		form.Set("include_links", "false")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	collection, _, err := client.Collection.Get(ctx, "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
	require.NoError(t, err)
	require.Equal(t, expectedCollection, collection)
}

func TestCollectionService_FromSubreddit(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/collection/collections.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/v1/collections/subreddit_collections", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("sr_fullname", "t5_2uquw1")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	collections, _, err := client.Collection.FromSubreddit(ctx, "t5_2uquw1")
	require.NoError(t, err)
	require.Equal(t, expectedCollections, collections)
}

func TestCollectionService_Create(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/collection/collection.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/v1/collections/create_collection", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("title", "Test Title")
		form.Set("sr_fullname", "t5_2uquw1")
		form.Set("display_layout", "TIMELINE")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	_, _, err = client.Collection.Create(ctx, nil)
	require.EqualError(t, err, "*CollectionCreateRequest: cannot be nil")

	collection, _, err := client.Collection.Create(ctx, &CollectionCreateRequest{
		Title:       "Test Title",
		SubredditID: "t5_2uquw1",
		Layout:      "TIMELINE",
	})
	require.NoError(t, err)
	require.Equal(t, expectedCollection, collection)
}

func TestCollectionService_Delete(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/collections/delete_collection", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("collection_id", "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Collection.Delete(ctx, "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
	require.NoError(t, err)
}

func TestCollectionService_AddPost(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/collections/add_post_to_collection", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("link_fullname", "t3_hs03f3")
		form.Set("collection_id", "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Collection.AddPost(ctx, "t3_hs03f3", "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
	require.NoError(t, err)
}

func TestCollectionService_RemovePost(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/collections/remove_post_in_collection", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("link_fullname", "t3_hs03f3")
		form.Set("collection_id", "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Collection.RemovePost(ctx, "t3_hs03f3", "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
	require.NoError(t, err)
}

func TestCollectionService_ReorderPosts(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/collections/reorder_collection", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("collection_id", "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
		form.Set("link_ids", "t3_hs0cyh,t3_hqrg8s,t3_hs03f3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Collection.ReorderPosts(ctx, "37f1e52d-7ec9-466b-b4cc-59e86e071ed7", "t3_hs0cyh", "t3_hqrg8s", "t3_hs03f3")
	require.NoError(t, err)
}

func TestCollectionService_UpdateTitle(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/collections/update_collection_title", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("collection_id", "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
		form.Set("title", "Test Title")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Collection.UpdateTitle(ctx, "37f1e52d-7ec9-466b-b4cc-59e86e071ed7", "Test Title")
	require.NoError(t, err)
}

func TestCollectionService_UpdateDescription(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/collections/update_collection_description", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("collection_id", "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
		form.Set("description", "Test Description")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Collection.UpdateDescription(ctx, "37f1e52d-7ec9-466b-b4cc-59e86e071ed7", "Test Description")
	require.NoError(t, err)
}

func TestCollectionService_UpdateLayoutTimeline(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/collections/update_collection_display_layout", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("collection_id", "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
		form.Set("display_layout", "TIMELINE")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Collection.UpdateLayoutTimeline(ctx, "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
	require.NoError(t, err)
}

func TestCollectionService_UpdateLayoutGallery(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/collections/update_collection_display_layout", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("collection_id", "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
		form.Set("display_layout", "GALLERY")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Collection.UpdateLayoutGallery(ctx, "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
	require.NoError(t, err)
}

func TestCollectionService_Follow(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/collections/follow_collection", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("collection_id", "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
		form.Set("follow", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Collection.Follow(ctx, "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
	require.NoError(t, err)
}

func TestCollectionService_Unfollow(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/collections/follow_collection", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("collection_id", "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
		form.Set("follow", "false")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Collection.Unfollow(ctx, "37f1e52d-7ec9-466b-b4cc-59e86e071ed7")
	require.NoError(t, err)
}
