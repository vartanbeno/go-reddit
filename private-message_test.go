package reddit

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrivateMessageService_ReadAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/read_all_messages", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusAccepted)
	})

	res, err := client.PrivateMessage.ReadAll(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, res.StatusCode)
}

func TestPrivateMessageService_Read(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/read_message", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test1,test2,test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)
	})

	_, err := client.PrivateMessage.Read(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	_, err = client.PrivateMessage.Read(ctx, "test1", "test2", "test3")
	require.NoError(t, err)
}

func TestPrivateMessageService_Unread(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unread_message", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test1,test2,test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)
	})

	_, err := client.PrivateMessage.Unread(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	_, err = client.PrivateMessage.Unread(ctx, "test1", "test2", "test3")
	require.NoError(t, err)
}

func TestPrivateMessageService_Block(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/block", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)
	})

	_, err := client.PrivateMessage.Block(ctx, "test")
	require.NoError(t, err)
}

func TestPrivateMessageService_Collapse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/collapse_message", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test1,test2,test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)
	})

	_, err := client.PrivateMessage.Collapse(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	_, err = client.PrivateMessage.Collapse(ctx, "test1", "test2", "test3")
	require.NoError(t, err)
}

func TestPrivateMessageService_Uncollapse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/uncollapse_message", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test1,test2,test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)
	})

	_, err := client.PrivateMessage.Uncollapse(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	_, err = client.PrivateMessage.Uncollapse(ctx, "test1", "test2", "test3")
	require.NoError(t, err)
}
